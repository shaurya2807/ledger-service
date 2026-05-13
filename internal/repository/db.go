package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shaurya2807/ledger-service/configs"
)

func NewPool(ctx context.Context, cfg *configs.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.DB.DSN())
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return pool, nil
}
