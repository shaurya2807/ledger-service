package model

import "time"

type Account struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}
