// @title           Transaction Service API
// @version         1.0
// @description     A service for managing cardholder accounts and transactions.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey ApiKeyAuth
// @in              header
// @name            X-API-Key

package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	_ "transaction-service/docs"
	"transaction-service/internal/config"
	"transaction-service/internal/db"
	"transaction-service/internal/router"
)

func main() {
	godotenv.Load() // load .env if present; ignored in production where env vars are set directly

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	database := db.Init()
	defer database.Close()

	log.Printf("Server running on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router.New(database, cfg)))
}
