package repository

import (
	"testing"

	"transaction-service/internal/testutil"
)

func TestTransactionRepository_Create(t *testing.T) {
	db := testutil.NewTestDB(t)
	accountRepo := NewAccountRepository(db)
	txRepo := NewTransactionRepository(db)

	acc, err := accountRepo.Create("12345678900")
	if err != nil {
		t.Fatal(err)
	}

	tx, err := txRepo.Create(acc.ID, 1, -50.0)
	if err != nil {
		t.Fatal(err)
	}
	if tx.ID == 0 {
		t.Error("expected non-zero transaction_id")
	}
	if tx.AccountID != acc.ID {
		t.Errorf("expected account_id %d, got %d", acc.ID, tx.AccountID)
	}
	if tx.Amount != -50.0 {
		t.Errorf("expected amount -50.0, got %f", tx.Amount)
	}
	if tx.EventDate.IsZero() {
		t.Error("expected non-zero event_date")
	}
}

func TestTransactionRepository_Create_InvalidAccount(t *testing.T) {
	db := testutil.NewTestDB(t)
	txRepo := NewTransactionRepository(db)

	if _, err := txRepo.Create(999, 1, -50.0); err == nil {
		t.Error("expected error for non-existent account_id, got nil")
	}
}

func TestTransactionRepository_OperationTypeExists(t *testing.T) {
	db := testutil.NewTestDB(t)
	txRepo := NewTransactionRepository(db)

	for _, id := range []int64{1, 2, 3, 4} {
		exists, err := txRepo.OperationTypeExists(id)
		if err != nil {
			t.Fatalf("operation_type_id %d: %v", id, err)
		}
		if !exists {
			t.Errorf("expected operation_type_id %d to exist", id)
		}
	}
}

func TestTransactionRepository_OperationTypeExists_Invalid(t *testing.T) {
	db := testutil.NewTestDB(t)
	txRepo := NewTransactionRepository(db)

	exists, err := txRepo.OperationTypeExists(99)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Error("expected operation_type_id 99 to not exist")
	}
}
