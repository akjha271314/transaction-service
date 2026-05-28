package service

import (
	"database/sql"

	"transaction-service/internal/models"
	"transaction-service/internal/repository"
)

type mockAccountRepo struct {
	createFn         func(documentNumber string, balance float64) (*models.Account, error)
	findByIDFn       func(id int64) (*models.Account, error)
	updateBalanceTxFn func(tx *sql.Tx, accountID int64, delta float64) error
}

func (m *mockAccountRepo) Create(documentNumber string, balance float64) (*models.Account, error) {
	return m.createFn(documentNumber, balance)
}

func (m *mockAccountRepo) FindByID(id int64) (*models.Account, error) {
	return m.findByIDFn(id)
}

func (m *mockAccountRepo) UpdateBalanceTx(tx *sql.Tx, accountID int64, delta float64) error {
	return m.updateBalanceTxFn(tx, accountID, delta)
}

type mockTransactionRepo struct {
	createTxFn            func(tx *sql.Tx, accountID, operationTypeID int64, amount float64) (*models.Transaction, error)
	operationTypeExistsFn func(id int64) (bool, error)
}

func (m *mockTransactionRepo) CreateTx(tx *sql.Tx, accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
	return m.createTxFn(tx, accountID, operationTypeID, amount)
}

func (m *mockTransactionRepo) OperationTypeExists(id int64) (bool, error) {
	return m.operationTypeExistsFn(id)
}

// mockTxRunner calls fn with nil — sufficient for unit tests that mock the DB methods.
type mockTxRunner struct{}

func (m *mockTxRunner) RunInTx(fn func(*sql.Tx) error) error {
	return fn(nil)
}

// helpers

func makeAccountRepo(acc *models.Account, err error) *mockAccountRepo {
	return &mockAccountRepo{
		findByIDFn:        func(id int64) (*models.Account, error) { return acc, err },
		updateBalanceTxFn: func(tx *sql.Tx, accountID int64, delta float64) error { return nil },
	}
}

func makeTxRepo(stored *float64) *mockTransactionRepo {
	return &mockTransactionRepo{
		operationTypeExistsFn: func(id int64) (bool, error) { return true, nil },
		createTxFn: func(tx *sql.Tx, accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
			if stored != nil {
				*stored = amount
			}
			return &models.Transaction{AccountID: accountID, OperationTypeID: operationTypeID, Amount: amount}, nil
		},
	}
}

func newSvc(txRepo *mockTransactionRepo, accountRepo *mockAccountRepo) TransactionService {
	return NewTransactionService(txRepo, accountRepo, &mockTxRunner{})
}

func insufficientBalanceRepo() *mockAccountRepo {
	return &mockAccountRepo{
		findByIDFn:        func(id int64) (*models.Account, error) { return &models.Account{ID: id}, nil },
		updateBalanceTxFn: func(tx *sql.Tx, accountID int64, delta float64) error { return repository.ErrInsufficientBalance },
	}
}
