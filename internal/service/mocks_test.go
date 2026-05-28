package service

import (
	"database/sql"

	"transaction-service/internal/models"
)

type mockAccountRepo struct {
	createFn   func(documentNumber string, creditLimit float64) (*models.Account, error)
	findByIDFn func(id int64) (*models.Account, error)
}

func (m *mockAccountRepo) Create(documentNumber string, creditLimit float64) (*models.Account, error) {
	return m.createFn(documentNumber, creditLimit)
}

func (m *mockAccountRepo) FindByID(id int64) (*models.Account, error) {
	return m.findByIDFn(id)
}

type mockTransactionRepo struct {
	createTxFn            func(tx *sql.Tx, accountID, operationTypeID int64, amount float64) (*models.Transaction, error)
	getBalanceTxFn        func(tx *sql.Tx, accountID int64) (float64, error)
	operationTypeExistsFn func(id int64) (bool, error)
}

func (m *mockTransactionRepo) CreateTx(tx *sql.Tx, accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
	return m.createTxFn(tx, accountID, operationTypeID, amount)
}

func (m *mockTransactionRepo) GetBalanceTx(tx *sql.Tx, accountID int64) (float64, error) {
	return m.getBalanceTxFn(tx, accountID)
}

func (m *mockTransactionRepo) OperationTypeExists(id int64) (bool, error) {
	return m.operationTypeExistsFn(id)
}

// mockTxRunner calls fn with nil — sufficient for unit tests that mock the DB methods.
type mockTxRunner struct{}

func (m *mockTxRunner) RunInTx(fn func(*sql.Tx) error) error {
	return fn(nil)
}
