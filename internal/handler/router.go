package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/shaurya2807/ledger-service/internal/service"
)

func NewRouter(
	accountSvc *service.AccountService,
	txSvc *service.TransactionService,
	logger *zap.Logger,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger(logger))

	r.GET("/health", Health)

	v1 := r.Group("/api/v1")
	{
		accounts := NewAccountHandler(accountSvc, logger)
		v1.POST("/accounts", accounts.CreateAccount)
		v1.GET("/accounts/:id", accounts.GetAccount)
		v1.GET("/accounts/:id/balance", accounts.GetBalance)

		txs := NewTransactionHandler(txSvc, logger)
		v1.POST("/transfers", txs.Transfer)
	}

	return r
}

func requestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}
