package service

import (
	"context"

	"github.com/shaurya2807/ledger-service/internal/model"
	"github.com/shaurya2807/ledger-service/internal/repository"
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context, req *model.CreateAccountRequest) (*model.Account, error) {
	return s.repo.Create(ctx, req)
}

func (s *AccountService) GetAccount(ctx context.Context, id string) (*model.Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AccountService) GetBalance(ctx context.Context, accountID string) (*model.BalanceResponse, error) {
	return s.repo.GetBalance(ctx, accountID)
}
