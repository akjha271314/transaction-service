package repository

import (
	"testing"

	"transaction-service/internal/testutil"
)

func TestAccountRepository_Create(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	acc, err := repo.Create("12345678900")
	if err != nil {
		t.Fatal(err)
	}
	if acc.ID == 0 {
		t.Error("expected non-zero account_id")
	}
	if acc.DocumentNumber != "12345678900" {
		t.Errorf("expected document_number 12345678900, got %s", acc.DocumentNumber)
	}
}

func TestAccountRepository_Create_DuplicateDocumentNumber(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	if _, err := repo.Create("12345678900"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Create("12345678900"); err == nil {
		t.Error("expected error on duplicate document_number, got nil")
	}
}

func TestAccountRepository_FindByID(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	created, err := repo.Create("12345678900")
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
	if found.DocumentNumber != created.DocumentNumber {
		t.Errorf("expected document_number %s, got %s", created.DocumentNumber, found.DocumentNumber)
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
