package repository

import (
	"testing"

	"transaction-service/internal/testutil"
)

func TestAccountRepository_Create(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	acc, err := repo.Create("12345678900", 500.0)
	if err != nil {
		t.Fatal(err)
	}
	if acc.ID == 0 {
		t.Error("expected non-zero account_id")
	}
	if acc.DocumentNumber != "12345678900" {
		t.Errorf("expected document_number 12345678900, got %s", acc.DocumentNumber)
	}
	if acc.Balance != 500.0 {
		t.Errorf("expected balance 500.0, got %f", acc.Balance)
	}
}

func TestAccountRepository_Create_DuplicateDocumentNumber(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	if _, err := repo.Create("12345678900", 0); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Create("12345678900", 0); err == nil {
		t.Error("expected error on duplicate document_number, got nil")
	}
}

func TestAccountRepository_FindByID(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	created, err := repo.Create("12345678900", 1000.0)
	if err != nil {
		t.Fatal(err)
	}

	found, err := repo.FindByID(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.ID != created.ID {
		t.Errorf("expected id %d, got %d", created.ID, found.ID)
	}
	if found.Balance != created.Balance {
		t.Errorf("expected balance %f, got %f", created.Balance, found.Balance)
	}
}

func TestAccountRepository_FindByID_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	_, err := repo.FindByID(999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAccountRepository_UpdateBalanceTx_Debit(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	acc, err := repo.Create("12345678900", 100.0)
	if err != nil {
		t.Fatal(err)
	}

	tx, _ := db.Begin()
	defer tx.Rollback()

	if err := repo.UpdateBalanceTx(tx, acc.ID, -40.0); err != nil {
		t.Fatal(err)
	}
	tx.Commit()

	updated, _ := repo.FindByID(acc.ID)
	if updated.Balance != 60.0 {
		t.Errorf("expected balance 60.0, got %f", updated.Balance)
	}
}

func TestAccountRepository_UpdateBalanceTx_InsufficientBalance(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	acc, err := repo.Create("12345678900", 30.0)
	if err != nil {
		t.Fatal(err)
	}

	tx, _ := db.Begin()
	defer tx.Rollback()

	err = repo.UpdateBalanceTx(tx, acc.ID, -50.0)
	if err != ErrInsufficientBalance {
		t.Errorf("expected ErrInsufficientBalance, got %v", err)
	}
}

func TestAccountRepository_UpdateBalanceTx_Credit(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	acc, err := repo.Create("12345678900", 0.0)
	if err != nil {
		t.Fatal(err)
	}

	tx, _ := db.Begin()
	defer tx.Rollback()

	if err := repo.UpdateBalanceTx(tx, acc.ID, 100.0); err != nil {
		t.Fatal(err)
	}
	tx.Commit()

	updated, _ := repo.FindByID(acc.ID)
	if updated.Balance != 100.0 {
		t.Errorf("expected balance 100.0, got %f", updated.Balance)
	}
}
