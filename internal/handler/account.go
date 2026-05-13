package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/shaurya2807/ledger-service/internal/model"
	"github.com/shaurya2807/ledger-service/internal/repository"
	"github.com/shaurya2807/ledger-service/internal/service"
)

type AccountHandler struct {
	svc    *service.AccountService
	logger *zap.Logger
}

func NewAccountHandler(svc *service.AccountService, logger *zap.Logger) *AccountHandler {
	return &AccountHandler{svc: svc, logger: logger}
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req model.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.svc.CreateAccount(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("create account failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("account created",
		zap.String("account_id", account.ID),
		zap.String("owner_id", account.OwnerID),
		zap.String("currency", account.Currency),
	)
	c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id: must be a UUID"})
		return
	}

	account, err := h.svc.GetAccount(c.Request.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}
	if err != nil {
		h.logger.Error("get account failed", zap.String("account_id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("account fetched", zap.String("account_id", account.ID))
	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) GetBalance(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id: must be a UUID"})
		return
	}

	resp, err := h.svc.GetBalance(c.Request.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}
	if err != nil {
		h.logger.Error("get balance failed", zap.String("account_id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("balance fetched",
		zap.String("account_id", resp.AccountID),
		zap.String("balance", resp.Balance),
	)
	c.JSON(http.StatusOK, resp)
}

// isValidUUID checks the standard 8-4-4-4-12 hex UUID format.
func isValidUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
	}
	return true
}
