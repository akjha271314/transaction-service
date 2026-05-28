package repository

import (
	"database/sql"
	"errors"
	"strings"

	"transaction-service/internal/models"
)

var ErrNotFound = errors.New("not found")
var ErrInsufficientBalance = errors.New("insufficient balance")
var ErrDuplicateAccount = errors.New("account with this document number already exists")

type AccountRepository interface {
	Create(documentNumber string, balance float64) (*models.Account, error)
	FindByID(id int64) (*models.Account, error)
	UpdateBalanceTx(tx *sql.Tx, accountID int64, delta float64) error
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(documentNumber string, balance float64) (*models.Account, error) {
	result, err := r.db.Exec(
		"INSERT INTO accounts (document_number, balance) VALUES (?, ?)",
		documentNumber, balance,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, ErrDuplicateAccount
		}
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &models.Account{ID: id, DocumentNumber: documentNumber, Balance: balance}, nil
}

func (r *accountRepository) FindByID(id int64) (*models.Account, error) {
	row := r.db.QueryRow(
		"SELECT account_id, document_number, balance FROM accounts WHERE account_id = ?", id,
	)
	var acc models.Account
	if err := row.Scan(&acc.ID, &acc.DocumentNumber, &acc.Balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &acc, nil
}

// UpdateBalanceTx atomically applies delta to the account balance within an existing
// transaction. The WHERE clause enforces balance >= 0, so a single affected row
// confirms both existence and sufficient funds — no separate read needed.
func (r *accountRepository) UpdateBalanceTx(tx *sql.Tx, accountID int64, delta float64) error {
	result, err := tx.Exec(
		"UPDATE accounts SET balance = balance + ? WHERE account_id = ? AND balance + ? >= 0",
		delta, accountID, delta,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrInsufficientBalance
	}
	return nil
}
