package service

import (
	"context"
	"fmt"

	"github.com/shaurya2807/ledger-service/internal/model"
	"github.com/shaurya2807/ledger-service/internal/repository"
)

type TransactionService struct {
	txRepo      *repository.TransactionRepository
	accountRepo *repository.AccountRepository
}

func NewTransactionService(
	txRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
) *TransactionService {
	return &TransactionService{txRepo: txRepo, accountRepo: accountRepo}
}

// Transfer executes a double-entry transfer. The second return value is true
// when the idempotency key was already seen; callers should respond with 409
// rather than 201 in that case.
func (s *TransactionService) Transfer(ctx context.Context, req *model.TransferRequest) (*model.Transaction, bool, error) {
	existing, err := s.txRepo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil {
		return nil, false, err
	}
	if existing != nil {
		return existing, true, nil
	}

	if req.FromAccountID == req.ToAccountID {
		return nil, false, ErrSameAccount
	}

	fromAccount, err := s.accountRepo.GetByID(ctx, req.FromAccountID)
	if err != nil {
		return nil, false, fmt.Errorf("from account: %w", err)
	}
	toAccount, err := s.accountRepo.GetByID(ctx, req.ToAccountID)
	if err != nil {
		return nil, false, fmt.Errorf("to account: %w", err)
	}

	if fromAccount.Currency != req.Currency {
		return nil, false, fmt.Errorf("%w: from account is %s, requested %s",
			ErrCurrencyMismatch, fromAccount.Currency, req.Currency)
	}
	if toAccount.Currency != req.Currency {
		return nil, false, fmt.Errorf("%w: to account is %s, requested %s",
			ErrCurrencyMismatch, toAccount.Currency, req.Currency)
	}

	tx, err := s.txRepo.Transfer(ctx, req)
	if err != nil {
		return nil, false, err
	}
	return tx, false, nil
}
