package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"transaction-service/internal/models"
	"transaction-service/internal/repository"
)

type mockAccountService struct {
	createFn func(documentNumber string) (*models.Account, error)
	getFn    func(id int64) (*models.Account, error)
}

func (m *mockAccountService) CreateAccount(documentNumber string) (*models.Account, error) {
	return m.createFn(documentNumber)
}

func (m *mockAccountService) GetAccount(id int64) (*models.Account, error) {
	return m.getFn(id)
}

func newAccountMux(svc *mockAccountService) *http.ServeMux {
	h := NewAccountHandler(svc)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /accounts", h.Create)
	mux.HandleFunc("GET /accounts/{accountId}", h.GetByID)
	return mux
}

func TestCreateAccount_Success(t *testing.T) {
	svc := &mockAccountService{
		createFn: func(documentNumber string) (*models.Account, error) {
			return &models.Account{ID: 1, DocumentNumber: documentNumber}, nil
		},
	}
	mux := newAccountMux(svc)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString(`{"document_number":"12345678900"}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var acc models.Account
	json.NewDecoder(w.Body).Decode(&acc)
	if acc.DocumentNumber != "12345678900" {
		t.Errorf("unexpected document_number: %s", acc.DocumentNumber)
	}
}

func TestCreateAccount_MissingDocumentNumber(t *testing.T) {
	svc := &mockAccountService{}
	mux := newAccountMux(svc)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString(`{}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetAccount_Success(t *testing.T) {
	svc := &mockAccountService{
		getFn: func(id int64) (*models.Account, error) {
			return &models.Account{ID: id, DocumentNumber: "12345678900"}, nil
		},
	}
	mux := newAccountMux(svc)

	req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var acc models.Account
	json.NewDecoder(w.Body).Decode(&acc)
	if acc.ID != 1 {
		t.Errorf("expected account_id 1, got %d", acc.ID)
	}
}

func TestGetAccount_NotFound(t *testing.T) {
	svc := &mockAccountService{
		getFn: func(id int64) (*models.Account, error) {
			return nil, repository.ErrNotFound
		},
	}
	mux := newAccountMux(svc)

	req := httptest.NewRequest(http.MethodGet, "/accounts/99", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetAccount_InvalidID(t *testing.T) {
	svc := &mockAccountService{}
	mux := newAccountMux(svc)

	req := httptest.NewRequest(http.MethodGet, "/accounts/abc", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}