package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

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
