package service

import "errors"

var (
	ErrSameAccount      = errors.New("from_account_id and to_account_id must be different")
	ErrCurrencyMismatch = errors.New("account currency does not match transfer currency")
)
