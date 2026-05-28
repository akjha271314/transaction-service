package router

import (
	"database/sql"
	"net/http"

	"transaction-service/internal/handler"
	"transaction-service/internal/repository"
	"transaction-service/internal/service"
)

func New(db *sql.DB) http.Handler {
	accountRepo := repository.NewAccountRepository(db)
	txRepo := repository.NewTransactionRepository(db)

	accountSvc := service.NewAccountService(accountRepo)
	txSvc := service.NewTransactionService(txRepo, accountRepo)

	accountHandler := handler.NewAccountHandler(accountSvc)
	txHandler := handler.NewTransactionHandler(txSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health)

	mux.HandleFunc("POST /accounts", accountHandler.Create)
	mux.HandleFunc("GET /accounts/{accountId}", accountHandler.GetByID)

	mux.HandleFunc("POST /transactions", txHandler.Create)

	return mux
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
