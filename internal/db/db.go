package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func Init() *sql.DB {
	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	if _, err := database.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatal("failed to enable foreign keys:", err)
	}
	if err := migrate(database); err != nil {
		log.Fatal("failed to run migrations:", err)
	}
	if err := seed(database); err != nil {
		log.Fatal("failed to seed database:", err)
	}
	return database
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			account_id      INTEGER PRIMARY KEY AUTOINCREMENT,
			document_number TEXT    NOT NULL UNIQUE,
			credit_limit    REAL    NOT NULL DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS operation_types (
			operation_type_id INTEGER PRIMARY KEY,
			description       TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS transactions (
			transaction_id    INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id        INTEGER NOT NULL REFERENCES accounts(account_id),
			operation_type_id INTEGER NOT NULL REFERENCES operation_types(operation_type_id),
			amount            REAL    NOT NULL,
			event_date        DATETIME NOT NULL
		);
	`)
	return err
}

func seed(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO operation_types (operation_type_id, description) VALUES
			(1, 'Normal Purchase'),
			(2, 'Purchase with installments'),
			(3, 'Withdrawal'),
			(4, 'Credit Voucher');
	`)
	return err
}
