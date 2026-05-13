package service

import (
	"context"
	"errors"
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

func (s *TransactionService) Transfer(ctx context.Context, req *model.TransferRequest) (*model.Transaction, error) {
	existing, err := s.txRepo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	if req.FromAccountID == req.ToAccountID {
		return nil, errors.New("from_account_id and to_account_id must differ")
	}

	fromAccount, err := s.accountRepo.GetByID(ctx, req.FromAccountID)
	if err != nil {
		return nil, fmt.Errorf("from account: %w", err)
	}
	toAccount, err := s.accountRepo.GetByID(ctx, req.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("to account: %w", err)
	}

	if fromAccount.Currency != req.Currency {
		return nil, fmt.Errorf("from account currency is %s, transfer currency is %s", fromAccount.Currency, req.Currency)
	}
	if toAccount.Currency != req.Currency {
		return nil, fmt.Errorf("to account currency is %s, transfer currency is %s", toAccount.Currency, req.Currency)
	}

	return s.txRepo.Transfer(ctx, req)
}
