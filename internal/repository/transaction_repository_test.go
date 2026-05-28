package repository

import (
	"testing"

	"transaction-service/internal/testutil"
)

func TestTransactionRepository_CreateTx(t *testing.T) {
	db := testutil.NewTestDB(t)
	accountRepo := NewAccountRepository(db)
	txRepo := NewTransactionRepository(db)

	acc, err := accountRepo.Create("12345678900", 500.0)
	if err != nil {
		t.Fatal(err)
	}

	sqlTx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlTx.Rollback()

	tx, err := txRepo.CreateTx(sqlTx, acc.ID, 1, -50.0)
	if err != nil {
		t.Fatal(err)
	}
	sqlTx.Commit()

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

func TestTransactionRepository_CreateTx_InvalidAccount(t *testing.T) {
	db := testutil.NewTestDB(t)
	txRepo := NewTransactionRepository(db)

	sqlTx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlTx.Rollback()

	if _, err := txRepo.CreateTx(sqlTx, 999, 1, -50.0); err == nil {
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
