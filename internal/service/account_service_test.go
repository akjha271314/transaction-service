package service

import (
	"errors"
	"testing"

	"transaction-service/internal/models"
	"transaction-service/internal/repository"
)

func TestCreateAccount(t *testing.T) {
	repo := &mockAccountRepo{
		createFn: func(documentNumber string, creditLimit float64) (*models.Account, error) {
			return &models.Account{ID: 1, DocumentNumber: documentNumber, CreditLimit: creditLimit}, nil
		},
	}
	svc := NewAccountService(repo)

	acc, err := svc.CreateAccount("12345678900", 500.0)
	if err != nil {
		t.Fatal(err)
	}
	if acc.DocumentNumber != "12345678900" {
		t.Errorf("expected document_number 12345678900, got %s", acc.DocumentNumber)
	}
	if acc.CreditLimit != 500.0 {
		t.Errorf("expected credit_limit 500.0, got %f", acc.CreditLimit)
	}
}

func TestGetAccount_Success(t *testing.T) {
	repo := &mockAccountRepo{
		findByIDFn: func(id int64) (*models.Account, error) {
			return &models.Account{ID: id, DocumentNumber: "12345678900"}, nil
		},
	}
	svc := NewAccountService(repo)

	acc, err := svc.GetAccount(1)
	if err != nil {
		t.Fatal(err)
	}
	if acc.ID != 1 {
		t.Errorf("expected id 1, got %d", acc.ID)
	}
}

func TestGetAccount_NotFound(t *testing.T) {
	repo := &mockAccountRepo{
		findByIDFn: func(id int64) (*models.Account, error) {
			return nil, repository.ErrNotFound
		},
	}
	svc := NewAccountService(repo)

	_, err := svc.GetAccount(99)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
