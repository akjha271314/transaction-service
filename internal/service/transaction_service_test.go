package service

import (
	"database/sql"
	"testing"

	"transaction-service/internal/models"
	"transaction-service/internal/repository"
)

func TestCreateTransaction_PurchaseAmount_IsNegative(t *testing.T) {
	var stored float64
	svc := newSvc(makeTxRepo(&stored), makeAccountRepo(&models.Account{ID: 1, Balance: 1000}, nil))

	if _, err := svc.CreateTransaction(1, 1, 50.0); err != nil {
		t.Fatal(err)
	}
	if stored != -50.0 {
		t.Errorf("expected -50.0, got %f", stored)
	}
}

func TestCreateTransaction_WithdrawalAmount_IsNegative(t *testing.T) {
	var stored float64
	svc := newSvc(makeTxRepo(&stored), makeAccountRepo(&models.Account{ID: 1, Balance: 1000}, nil))

	if _, err := svc.CreateTransaction(1, 3, 20.0); err != nil {
		t.Fatal(err)
	}
	if stored != -20.0 {
		t.Errorf("expected -20.0, got %f", stored)
	}
}

func TestCreateTransaction_CreditVoucher_IsPositive(t *testing.T) {
	var stored float64
	svc := newSvc(makeTxRepo(&stored), makeAccountRepo(&models.Account{ID: 1, Balance: 0}, nil))

	if _, err := svc.CreateTransaction(1, 4, 60.0); err != nil {
		t.Fatal(err)
	}
	if stored != 60.0 {
		t.Errorf("expected 60.0, got %f", stored)
	}
}

func TestCreateTransaction_NegativeInputAlwaysNormalized(t *testing.T) {
	var stored float64
	svc := newSvc(makeTxRepo(&stored), makeAccountRepo(&models.Account{ID: 1, Balance: 0}, nil))

	// Caller sends -60 for credit voucher — should be stored as +60
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
		findOperationTypeFn: func(id int64) (*models.OperationType, error) { return nil, repository.ErrNotFound },
	}
	svc := newSvc(txRepo, makeAccountRepo(&models.Account{ID: 1}, nil))

	_, err := svc.CreateTransaction(1, 99, 50.0)
	if err != ErrInvalidOperationType {
		t.Errorf("expected ErrInvalidOperationType, got %v", err)
	}
}

func TestCreateTransaction_InsufficientBalance(t *testing.T) {
	txRepo := &mockTransactionRepo{
		findOperationTypeFn: func(id int64) (*models.OperationType, error) {
			return &models.OperationType{ID: id, IsCredit: false}, nil
		},
		createTxFn: func(tx *sql.Tx, a, b int64, c float64) (*models.Transaction, error) { return nil, nil },
	}
	svc := newSvc(txRepo, insufficientBalanceRepo())

	_, err := svc.CreateTransaction(1, 1, 150.0)
	if err != ErrInsufficientBalance {
		t.Errorf("expected ErrInsufficientBalance, got %v", err)
	}
}
