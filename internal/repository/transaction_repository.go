package repository

import (
	"database/sql"
	"time"

	"transaction-service/internal/models"
)

type TransactionRepository interface {
	CreateTx(tx *sql.Tx, accountID, operationTypeID int64, amount float64) (*models.Transaction, error)
	GetBalanceTx(tx *sql.Tx, accountID int64) (float64, error)
	OperationTypeExists(id int64) (bool, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTx(tx *sql.Tx, accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
	eventDate := time.Now().UTC()
	result, err := tx.Exec(
		"INSERT INTO transactions (account_id, operation_type_id, amount, event_date) VALUES (?, ?, ?, ?)",
		accountID, operationTypeID, amount, eventDate,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &models.Transaction{
		ID:              id,
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          amount,
		EventDate:       eventDate,
	}, nil
}

func (r *transactionRepository) GetBalanceTx(tx *sql.Tx, accountID int64) (float64, error) {
	var balance float64
	err := tx.QueryRow(
		"SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE account_id = ?", accountID,
	).Scan(&balance)
	return balance, err
}

func (r *transactionRepository) OperationTypeExists(id int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(1) FROM operation_types WHERE operation_type_id = ?", id,
	).Scan(&count)
	return count > 0, err
}
