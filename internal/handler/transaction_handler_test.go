package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"transaction-service/internal/models"
	"transaction-service/internal/repository"
	"transaction-service/internal/service"
)

type mockTransactionService struct {
	createFn func(accountID, operationTypeID int64, amount float64) (*models.Transaction, error)
}

func (m *mockTransactionService) CreateTransaction(accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
	return m.createFn(accountID, operationTypeID, amount)
}

func newTxMux(svc *mockTransactionService) *http.ServeMux {
	h := NewTransactionHandler(svc)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /transactions", h.Create)
	return mux
}

func TestCreateTransaction_Success(t *testing.T) {
	svc := &mockTransactionService{
		createFn: func(accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
			return &models.Transaction{
				ID:              1,
				AccountID:       accountID,
				OperationTypeID: operationTypeID,
				Amount:          amount,
				EventDate:       time.Now(),
			}, nil
		},
	}
	mux := newTxMux(svc)

	body := `{"account_id":1,"operation_type_id":4,"amount":123.45}`
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var tx models.Transaction
	json.NewDecoder(w.Body).Decode(&tx)
	if tx.Amount != 123.45 {
		t.Errorf("expected amount 123.45, got %f", tx.Amount)
	}
}

func TestCreateTransaction_InvalidAccount(t *testing.T) {
	svc := &mockTransactionService{
		createFn: func(accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
			return nil, service.ErrInvalidAccount
		},
	}
	mux := newTxMux(svc)

	body := `{"account_id":99,"operation_type_id":1,"amount":50.0}`
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestCreateTransaction_InvalidOperationType(t *testing.T) {
	svc := &mockTransactionService{
		createFn: func(accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
			return nil, service.ErrInvalidOperationType
		},
	}
	mux := newTxMux(svc)

	body := `{"account_id":1,"operation_type_id":99,"amount":50.0}`
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestCreateTransaction_InsufficientBalance(t *testing.T) {
	svc := &mockTransactionService{
		createFn: func(accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
			return nil, repository.ErrInsufficientBalance
		},
	}
	mux := newTxMux(svc)

	body := `{"account_id":1,"operation_type_id":1,"amount":9999}`
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestCreateTransaction_InvalidBody(t *testing.T) {
	svc := &mockTransactionService{}
	mux := newTxMux(svc)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString("not json"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}