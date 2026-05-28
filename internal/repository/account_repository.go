package repository

import (
	"database/sql"
	"errors"

	"transaction-service/internal/models"
)

var ErrNotFound = errors.New("not found")

type AccountRepository interface {
	Create(documentNumber string) (*models.Account, error)
	FindByID(id int64) (*models.Account, error)
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(documentNumber string) (*models.Account, error) {
	result, err := r.db.Exec(
		"INSERT INTO accounts (document_number) VALUES (?)",
		documentNumber,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &models.Account{ID: id, DocumentNumber: documentNumber}, nil
}

func (r *accountRepository) FindByID(id int64) (*models.Account, error) {
	row := r.db.QueryRow(
		"SELECT account_id, document_number FROM accounts WHERE account_id = ?", id,
	)
	var acc models.Account
	if err := row.Scan(&acc.ID, &acc.DocumentNumber); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &acc, nil
}
