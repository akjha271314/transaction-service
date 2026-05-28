package service

import (
	"database/sql"
	"testing"

	"transaction-service/internal/models"
	"transaction-service/internal/repository"
)

func makeAccountRepo(acc *models.Account, err error) *mockAccountRepo {
	return &mockAccountRepo{
		findByIDFn: func(id int64) (*models.Account, error) { return acc, err },
	}
}

func makeTxRepo(stored *float64) *mockTransactionRepo {
	return &mockTransactionRepo{
		operationTypeExistsFn: func(id int64) (bool, error) { return true, nil },
		getBalanceTxFn:        func(tx *sql.Tx, accountID int64) (float64, error) { return 0, nil },
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

func TestCreateTransaction_PurchaseAmount_IsNegative(t *testing.T) {
	var stored float64
	svc := newSvc(makeTxRepo(&stored), makeAccountRepo(&models.Account{ID: 1, CreditLimit: 1000}, nil))

	if _, err := svc.CreateTransaction(1, 1, 50.0); err != nil {
		t.Fatal(err)
	}
	if stored != -50.0 {
		t.Errorf("expected -50.0, got %f", stored)
	}
}

func TestCreateTransaction_WithdrawalAmount_IsNegative(t *testing.T) {
	var stored float64
	svc := newSvc(makeTxRepo(&stored), makeAccountRepo(&models.Account{ID: 1, CreditLimit: 1000}, nil))

	if _, err := svc.CreateTransaction(1, 3, 20.0); err != nil {
		t.Fatal(err)
	}
	if stored != -20.0 {
		t.Errorf("expected -20.0, got %f", stored)
	}
}

func TestCreateTransaction_CreditVoucher_IsPositive(t *testing.T) {
	var stored float64
	svc := newSvc(makeTxRepo(&stored), makeAccountRepo(&models.Account{ID: 1, CreditLimit: 1000}, nil))

	if _, err := svc.CreateTransaction(1, 4, 60.0); err != nil {
		t.Fatal(err)
	}
	if stored != 60.0 {
		t.Errorf("expected 60.0, got %f", stored)
	}
}

func TestCreateTransaction_NegativeInputAlwaysNormalized(t *testing.T) {
	var stored float64
	svc := newSvc(makeTxRepo(&stored), makeAccountRepo(&models.Account{ID: 1, CreditLimit: 1000}, nil))

	if _, err := svc.CreateTransaction(1, 4, -60.0); err != nil {
		t.Fatal(err)
	}
	if stored != 60.0 {
		t.Errorf("expected 60.0, got %f", stored)
	}
}

func TestCreateTransaction_InvalidAccount(t *testing.T) {
	svc := newSvc(makeTxRepo(nil), makeAccountRepo(nil, repository.ErrNotFound))

	_, err := svc.CreateTransaction(99, 1, 50.0)
	if err != ErrInvalidAccount {
		t.Errorf("expected ErrInvalidAccount, got %v", err)
	}
}

func TestCreateTransaction_InvalidOperationType(t *testing.T) {
	txRepo := &mockTransactionRepo{
		operationTypeExistsFn: func(id int64) (bool, error) { return false, nil },
	}
	svc := newSvc(txRepo, makeAccountRepo(&models.Account{ID: 1}, nil))

	_, err := svc.CreateTransaction(1, 99, 50.0)
	if err != ErrInvalidOperationType {
		t.Errorf("expected ErrInvalidOperationType, got %v", err)
	}
}

func TestCreateTransaction_ExceedsCreditLimit(t *testing.T) {
	txRepo := &mockTransactionRepo{
		operationTypeExistsFn: func(id int64) (bool, error) { return true, nil },
		getBalanceTxFn:        func(tx *sql.Tx, accountID int64) (float64, error) { return -400.0, nil },
		createTxFn:            func(tx *sql.Tx, accountID, operationTypeID int64, amount float64) (*models.Transaction, error) { return nil, nil },
	}
	// credit_limit = 500, current balance = -400, new purchase = -150 → total -550 < -500
	svc := newSvc(txRepo, makeAccountRepo(&models.Account{ID: 1, CreditLimit: 500}, nil))

	_, err := svc.CreateTransaction(1, 1, 150.0)
	if err != ErrInsufficientCredit {
		t.Errorf("expected ErrInsufficientCredit, got %v", err)
	}
}

func TestCreateTransaction_WithinCreditLimit(t *testing.T) {
	txRepo := &mockTransactionRepo{
		operationTypeExistsFn: func(id int64) (bool, error) { return true, nil },
		getBalanceTxFn:        func(tx *sql.Tx, accountID int64) (float64, error) { return -400.0, nil },
		createTxFn: func(tx *sql.Tx, accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
			return &models.Transaction{Amount: amount}, nil
		},
	}
	// credit_limit = 500, current balance = -400, new purchase = -50 → total -450 >= -500 ✓
	svc := newSvc(txRepo, makeAccountRepo(&models.Account{ID: 1, CreditLimit: 500}, nil))

	_, err := svc.CreateTransaction(1, 1, 50.0)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
