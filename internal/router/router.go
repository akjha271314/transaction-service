package router

import (
	"database/sql"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"transaction-service/internal/config"
	"transaction-service/internal/handler"
	"transaction-service/internal/repository"
	"transaction-service/internal/service"
)

func New(db *sql.DB, cfg *config.Config) http.Handler {
	accountRepo := repository.NewAccountRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	txRunner := repository.NewTxRunner(db)

	accountSvc := service.NewAccountService(accountRepo)
	txSvc := service.NewTransactionService(txRepo, accountRepo, txRunner)

	accountHandler := handler.NewAccountHandler(accountSvc)
	txHandler := handler.NewTransactionHandler(txSvc)

	protected := http.NewServeMux()
	protected.HandleFunc("POST /accounts", accountHandler.Create)
	protected.HandleFunc("GET /accounts/{accountId}", accountHandler.GetByID)
	protected.HandleFunc("POST /transactions", txHandler.Create)

	top := http.NewServeMux()
	top.HandleFunc("GET /health", health)
	top.Handle("/swagger/", httpSwagger.WrapHandler)
	top.Handle("/", apiKeyMiddleware(cfg.APIKey, protected))

	return top
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
