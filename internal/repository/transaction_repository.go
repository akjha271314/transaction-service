package repository

import (
	"database/sql"
	"errors"
	"time"

	"transaction-service/internal/models"
)

type TransactionRepository interface {
	CreateTx(tx *sql.Tx, accountID, operationTypeID int64, amount float64) (*models.Transaction, error)
	FindOperationType(id int64) (*models.OperationType, error)
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

func (r *transactionRepository) FindOperationType(id int64) (*models.OperationType, error) {
	row := r.db.QueryRow(
		"SELECT operation_type_id, description, is_credit FROM operation_types WHERE operation_type_id = ?", id,
	)
	var op models.OperationType
	if err := row.Scan(&op.ID, &op.Description, &op.IsCredit); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &op, nil
}
