package main

import (
	"log"
	"net/http"

	"transaction-service/internal/db"
	"transaction-service/internal/router"
)

func main() {
	database := db.Init()
	defer database.Close()

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router.New(database)))
}
