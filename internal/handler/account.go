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
		h.logger.Error("create account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	id := c.Param("id")

	account, err := h.svc.GetAccount(c.Request.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}
	if err != nil {
		h.logger.Error("get account", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) GetBalance(c *gin.Context) {
	id := c.Param("id")

	balance, err := h.svc.GetBalance(c.Request.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}
	if err != nil {
		h.logger.Error("get balance", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, balance)
}
