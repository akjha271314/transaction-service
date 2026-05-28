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

func TestTransactionRepository_FindOperationType(t *testing.T) {
	db := testutil.NewTestDB(t)
	txRepo := NewTransactionRepository(db)

	cases := []struct {
		id       int64
		isCredit bool
	}{
		{1, false},
		{2, false},
		{3, false},
		{4, true},
	}
	for _, tc := range cases {
		op, err := txRepo.FindOperationType(tc.id)
		if err != nil {
			t.Fatalf("operation_type_id %d: %v", tc.id, err)
		}
		if op.IsCredit != tc.isCredit {
			t.Errorf("operation_type_id %d: expected is_credit=%v, got %v", tc.id, tc.isCredit, op.IsCredit)
		}
	}
}

func TestTransactionRepository_FindOperationType_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	txRepo := NewTransactionRepository(db)

	_, err := txRepo.FindOperationType(99)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
