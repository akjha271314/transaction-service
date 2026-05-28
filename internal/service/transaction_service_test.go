package service

import (
	"testing"

	"transaction-service/internal/models"
	"transaction-service/internal/repository"
)

func makeAccountRepo(acc *models.Account, err error) *mockAccountRepo {
	return &mockAccountRepo{
		findByIDFn: func(id int64) (*models.Account, error) { return acc, err },
	}
}

func makeTxRepo(amount *float64) *mockTransactionRepo {
	return &mockTransactionRepo{
		operationTypeExistsFn: func(id int64) (bool, error) { return true, nil },
		createFn: func(accountID, operationTypeID int64, a float64) (*models.Transaction, error) {
			if amount != nil {
				*amount = a
			}
			return &models.Transaction{AccountID: accountID, OperationTypeID: operationTypeID, Amount: a}, nil
		},
	}
}

func TestCreateTransaction_PurchaseAmount_IsNegative(t *testing.T) {
	var stored float64
	svc := NewTransactionService(makeTxRepo(&stored), makeAccountRepo(&models.Account{ID: 1}, nil))

	if _, err := svc.CreateTransaction(1, 1, 50.0); err != nil {
		t.Fatal(err)
	}
	if stored != -50.0 {
		t.Errorf("expected -50.0, got %f", stored)
	}
}

func TestCreateTransaction_WithdrawalAmount_IsNegative(t *testing.T) {
	var stored float64
	svc := NewTransactionService(makeTxRepo(&stored), makeAccountRepo(&models.Account{ID: 1}, nil))

	if _, err := svc.CreateTransaction(1, 3, 20.0); err != nil {
		t.Fatal(err)
	}
	if stored != -20.0 {
		t.Errorf("expected -20.0, got %f", stored)
	}
}

func TestCreateTransaction_CreditVoucher_IsPositive(t *testing.T) {
	var stored float64
	svc := NewTransactionService(makeTxRepo(&stored), makeAccountRepo(&models.Account{ID: 1}, nil))

	if _, err := svc.CreateTransaction(1, 4, 60.0); err != nil {
		t.Fatal(err)
	}
	if stored != 60.0 {
		t.Errorf("expected 60.0, got %f", stored)
	}
}

func TestCreateTransaction_NegativeInputAlwaysNormalized(t *testing.T) {
	var stored float64
	svc := NewTransactionService(makeTxRepo(&stored), makeAccountRepo(&models.Account{ID: 1}, nil))

	// Caller sends -60.0 for credit voucher — should be stored as +60.0
	if _, err := svc.CreateTransaction(1, 4, -60.0); err != nil {
		t.Fatal(err)
	}
	if stored != 60.0 {
		t.Errorf("expected 60.0, got %f", stored)
	}
}

func TestCreateTransaction_InvalidAccount(t *testing.T) {
	svc := NewTransactionService(makeTxRepo(nil), makeAccountRepo(nil, repository.ErrNotFound))

	_, err := svc.CreateTransaction(99, 1, 50.0)
	if err != ErrInvalidAccount {
		t.Errorf("expected ErrInvalidAccount, got %v", err)
	}
}

func TestCreateTransaction_InvalidOperationType(t *testing.T) {
	txRepo := &mockTransactionRepo{
		operationTypeExistsFn: func(id int64) (bool, error) { return false, nil },
	}
	svc := NewTransactionService(txRepo, makeAccountRepo(&models.Account{ID: 1}, nil))

	_, err := svc.CreateTransaction(1, 99, 50.0)
	if err != ErrInvalidOperationType {
		t.Errorf("expected ErrInvalidOperationType, got %v", err)
	}
}