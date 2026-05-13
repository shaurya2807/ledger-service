package model

type CreateAccountRequest struct {
	OwnerID  string `json:"owner_id"  binding:"required"`
	Currency string `json:"currency"  binding:"required,len=3"`
}

type SeedRequest struct {
	Amount   float64 `json:"amount"   binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required,len=3"`
}

type TransferRequest struct {
	FromAccountID  string  `json:"from_account_id"  binding:"required,uuid"`
	ToAccountID    string  `json:"to_account_id"    binding:"required,uuid"`
	Amount         float64 `json:"amount"           binding:"required,gt=0"`
	Currency       string  `json:"currency"         binding:"required,len=3"`
	IdempotencyKey string  `json:"idempotency_key"  binding:"required"`
}
