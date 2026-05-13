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

type TransactionHandler struct {
	svc    *service.TransactionService
	logger *zap.Logger
}

func NewTransactionHandler(svc *service.TransactionService, logger *zap.Logger) *TransactionHandler {
	return &TransactionHandler{svc: svc, logger: logger}
}

func (h *TransactionHandler) Transfer(c *gin.Context) {
	var req model.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, isDuplicate, err := h.svc.Transfer(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInsufficientFunds):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "insufficient funds"})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		case errors.Is(err, service.ErrSameAccount):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrCurrencyMismatch):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			h.logger.Error("transfer failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	if isDuplicate {
		c.JSON(http.StatusConflict, tx)
		return
	}

	h.logger.Info("transfer completed",
		zap.String("transaction_id", tx.ID),
		zap.String("from_account_id", tx.FromAccountID),
		zap.String("to_account_id", tx.ToAccountID),
		zap.String("amount", tx.Amount),
	)
	c.JSON(http.StatusCreated, tx)
}
