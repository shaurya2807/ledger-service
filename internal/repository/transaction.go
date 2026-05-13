package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shaurya2807/ledger-service/internal/model"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetByIdempotencyKey(ctx context.Context, key string) (*model.Transaction, error) {
	var t model.Transaction
	err := r.db.QueryRow(ctx,
		`SELECT id, idempotency_key, from_account_id, to_account_id,
		        amount::TEXT, currency, status, created_at
		 FROM transactions WHERE idempotency_key = $1`,
		key,
	).Scan(&t.ID, &t.IdempotencyKey, &t.FromAccountID, &t.ToAccountID,
		&t.Amount, &t.Currency, &t.Status, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get transaction by idempotency key: %w", err)
	}
	return &t, nil
}

func (r *TransactionRepository) Transfer(ctx context.Context, req *model.TransferRequest) (*model.Transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// Lock all existing entry rows for from_account so that concurrent transfers
	// against the same account queue behind this transaction rather than racing
	// the balance check below.
	if _, err = tx.Exec(ctx,
		`SELECT id FROM entries WHERE account_id = $1 FOR UPDATE`,
		req.FromAccountID,
	); err != nil {
		return nil, fmt.Errorf("lock entries: %w", err)
	}

	// Compute current balance from the locked snapshot.
	var balance float64
	if err = tx.QueryRow(ctx,
		`SELECT COALESCE(
			SUM(CASE WHEN direction = 'credit' THEN amount ELSE -amount END),
			0
		)::FLOAT8
		 FROM entries WHERE account_id = $1`,
		req.FromAccountID,
	).Scan(&balance); err != nil {
		return nil, fmt.Errorf("compute balance: %w", err)
	}

	if balance < req.Amount {
		return nil, ErrInsufficientFunds
	}

	// Insert the transaction record first to obtain its generated ID.
	var t model.Transaction
	if err = tx.QueryRow(ctx,
		`INSERT INTO transactions (idempotency_key, from_account_id, to_account_id, amount, currency)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, idempotency_key, from_account_id, to_account_id,
		           amount::TEXT, currency, status, created_at`,
		req.IdempotencyKey, req.FromAccountID, req.ToAccountID, req.Amount, req.Currency,
	).Scan(&t.ID, &t.IdempotencyKey, &t.FromAccountID, &t.ToAccountID,
		&t.Amount, &t.Currency, &t.Status, &t.CreatedAt); err != nil {
		return nil, fmt.Errorf("insert transaction: %w", err)
	}

	// Debit from_account.
	if _, err = tx.Exec(ctx,
		`INSERT INTO entries (account_id, transaction_id, amount, direction)
		 VALUES ($1, $2, $3, 'debit')`,
		req.FromAccountID, t.ID, req.Amount,
	); err != nil {
		return nil, fmt.Errorf("insert debit entry: %w", err)
	}

	// Credit to_account.
	if _, err = tx.Exec(ctx,
		`INSERT INTO entries (account_id, transaction_id, amount, direction)
		 VALUES ($1, $2, $3, 'credit')`,
		req.ToAccountID, t.ID, req.Amount,
	); err != nil {
		return nil, fmt.Errorf("insert credit entry: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	return &t, nil
}
