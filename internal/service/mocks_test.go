package service

import "transaction-service/internal/models"

type mockAccountRepo struct {
	createFn   func(documentNumber string) (*models.Account, error)
	findByIDFn func(id int64) (*models.Account, error)
}

func (m *mockAccountRepo) Create(documentNumber string) (*models.Account, error) {
	return m.createFn(documentNumber)
}

func (m *mockAccountRepo) FindByID(id int64) (*models.Account, error) {
	return m.findByIDFn(id)
}

type mockTransactionRepo struct {
	createFn              func(accountID, operationTypeID int64, amount float64) (*models.Transaction, error)
	operationTypeExistsFn func(id int64) (bool, error)
}

func (m *mockTransactionRepo) Create(accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
	return m.createFn(accountID, operationTypeID, amount)
}

func (m *mockTransactionRepo) OperationTypeExists(id int64) (bool, error) {
	return m.operationTypeExistsFn(id)
}