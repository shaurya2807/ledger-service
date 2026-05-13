package model

import "time"

type Entry struct {
	ID            string    `json:"id"`
	AccountID     string    `json:"account_id"`
	TransactionID string    `json:"transaction_id"`
	Amount        string    `json:"amount"`
	Direction     string    `json:"direction"`
	CreatedAt     time.Time `json:"created_at"`
}
