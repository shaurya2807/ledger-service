package model

type BalanceResponse struct {
	AccountID  string `json:"account_id"`
	Currency   string `json:"currency"`
	Balance    string `json:"balance"`
	EntryCount int64  `json:"entry_count"`
}
