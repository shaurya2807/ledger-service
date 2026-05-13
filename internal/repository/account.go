package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shaurya2807/ledger-service/internal/model"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, req *model.CreateAccountRequest) (*model.Account, error) {
	var a model.Account
	err := r.db.QueryRow(ctx,
		`INSERT INTO accounts (owner_id, currency)
		 VALUES ($1, $2)
		 RETURNING id, owner_id, currency, created_at`,
		req.OwnerID, req.Currency,
	).Scan(&a.ID, &a.OwnerID, &a.Currency, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}
	return &a, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id string) (*model.Account, error) {
	var a model.Account
	err := r.db.QueryRow(ctx,
		`SELECT id, owner_id, currency, created_at FROM accounts WHERE id = $1`,
		id,
	).Scan(&a.ID, &a.OwnerID, &a.Currency, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}
	return &a, nil
}

func (r *AccountRepository) Seed(ctx context.Context, accountID string, req *model.SeedRequest) (*model.Entry, error) {
	var e model.Entry
	err := r.db.QueryRow(ctx,
		`INSERT INTO entries (account_id, transaction_id, amount, direction)
		 VALUES ($1, gen_random_uuid(), $2, 'credit')
		 RETURNING id, account_id, transaction_id, amount::TEXT, direction, created_at`,
		accountID, req.Amount,
	).Scan(&e.ID, &e.AccountID, &e.TransactionID, &e.Amount, &e.Direction, &e.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("seed account: %w", err)
	}
	return &e, nil
}

// GetBalance returns the computed balance and entry count for an account.
// It queries entries directly; callers must verify the account exists beforehand.
func (r *AccountRepository) GetBalance(ctx context.Context, accountID string) (balance string, entryCount int64, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT
			COUNT(*) AS entry_count,
			COALESCE(
				SUM(CASE WHEN direction = 'credit' THEN amount ELSE -amount END),
				0
			)::TEXT AS balance
		 FROM entries
		 WHERE account_id = $1`,
		accountID,
	).Scan(&entryCount, &balance)
	if err != nil {
		return "", 0, fmt.Errorf("get balance: %w", err)
	}
	return balance, entryCount, nil
}
