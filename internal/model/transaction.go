package model

import "time"

type Transaction struct {
	ID             string    `json:"id"`
	IdempotencyKey string    `json:"idempotency_key"`
	FromAccountID  string    `json:"from_account_id"`
	ToAccountID    string    `json:"to_account_id"`
	Amount         string    `json:"amount"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}
