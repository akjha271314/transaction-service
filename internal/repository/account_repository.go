package repository

import (
	"database/sql"
	"errors"

	"transaction-service/internal/models"
)

var ErrNotFound = errors.New("not found")

type AccountRepository interface {
	Create(documentNumber string, creditLimit float64) (*models.Account, error)
	FindByID(id int64) (*models.Account, error)
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(documentNumber string, creditLimit float64) (*models.Account, error) {
	result, err := r.db.Exec(
		"INSERT INTO accounts (document_number, credit_limit) VALUES (?, ?)",
		documentNumber, creditLimit,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &models.Account{ID: id, DocumentNumber: documentNumber, CreditLimit: creditLimit}, nil
}

func (r *accountRepository) FindByID(id int64) (*models.Account, error) {
	row := r.db.QueryRow(
		"SELECT account_id, document_number, credit_limit FROM accounts WHERE account_id = ?", id,
	)
	var acc models.Account
	if err := row.Scan(&acc.ID, &acc.DocumentNumber, &acc.CreditLimit); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &acc, nil
}
