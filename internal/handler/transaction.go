package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/shaurya2807/ledger-service/internal/model"
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

	tx, err := h.svc.Transfer(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("transfer", zap.Error(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}
