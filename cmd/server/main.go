package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"github.com/shaurya2807/ledger-service/configs"
	"github.com/shaurya2807/ledger-service/internal/handler"
	"github.com/shaurya2807/ledger-service/internal/repository"
	"github.com/shaurya2807/ledger-service/internal/service"
	"github.com/shaurya2807/ledger-service/pkg/logger"
)

func main() {
	_ = godotenv.Load()

	cfg, err := configs.Load()
	if err != nil {
		panic("load config: " + err.Error())
	}

	log, err := logger.New(cfg.AppEnv)
	if err != nil {
		panic("init logger: " + err.Error())
	}
	defer log.Sync() //nolint:errcheck

	ctx := context.Background()
	db, err := repository.NewPool(ctx, cfg)
	if err != nil {
		log.Fatal("connect to database", zap.Error(err))
	}
	defer db.Close()

	accountRepo := repository.NewAccountRepository(db)
	txRepo := repository.NewTransactionRepository(db)

	accountSvc := service.NewAccountService(accountRepo)
	txSvc := service.NewTransactionService(txRepo, accountRepo)

	router := handler.NewRouter(accountSvc, txSvc, log)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server starting", zap.String("port", cfg.ServerPort), zap.String("env", cfg.AppEnv))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen and serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutdown signal received")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("graceful shutdown failed", zap.Error(err))
	}
	log.Info("server stopped")
}
