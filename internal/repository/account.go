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

func (r *AccountRepository) GetBalance(ctx context.Context, accountID string) (*model.BalanceResponse, error) {
	var resp model.BalanceResponse
	err := r.db.QueryRow(ctx,
		`SELECT
			a.id,
			a.currency,
			COALESCE(SUM(CASE WHEN e.direction = 'credit' THEN e.amount ELSE -e.amount END), 0)::TEXT,
			COUNT(e.id)
		 FROM accounts a
		 LEFT JOIN entries e ON e.account_id = a.id
		 WHERE a.id = $1
		 GROUP BY a.id, a.currency`,
		accountID,
	).Scan(&resp.AccountID, &resp.Currency, &resp.Balance, &resp.EntryCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}
	return &resp, nil
}
