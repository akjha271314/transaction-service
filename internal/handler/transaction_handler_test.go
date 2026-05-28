package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"transaction-service/internal/models"
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

func TestTransactionHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mockCreate func(int64, int64, float64) (*models.Transaction, error)
		wantStatus int
		wantAmount float64
	}{
		{
			name: "purchase success",
			body: `{"account_id":1,"operation_type_id":1,"amount":50}`,
			mockCreate: func(accountID, opTypeID int64, amount float64) (*models.Transaction, error) {
				return &models.Transaction{ID: 1, AccountID: accountID, OperationTypeID: opTypeID, Amount: amount, EventDate: time.Now()}, nil
			},
			wantStatus: http.StatusCreated,
			wantAmount: 50,
		},
		{
			name: "credit voucher success",
			body: `{"account_id":1,"operation_type_id":4,"amount":123.45}`,
			mockCreate: func(accountID, opTypeID int64, amount float64) (*models.Transaction, error) {
				return &models.Transaction{ID: 2, AccountID: accountID, OperationTypeID: opTypeID, Amount: amount, EventDate: time.Now()}, nil
			},
			wantStatus: http.StatusCreated,
			wantAmount: 123.45,
		},
		{
			name: "invalid account",
			body: `{"account_id":99,"operation_type_id":1,"amount":50}`,
			mockCreate: func(a, b int64, c float64) (*models.Transaction, error) {
				return nil, service.ErrInvalidAccount
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid operation type",
			body: `{"account_id":1,"operation_type_id":99,"amount":50}`,
			mockCreate: func(a, b int64, c float64) (*models.Transaction, error) {
				return nil, service.ErrInvalidOperationType
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "insufficient balance",
			body: `{"account_id":1,"operation_type_id":1,"amount":9999}`,
			mockCreate: func(a, b int64, c float64) (*models.Transaction, error) {
				return nil, service.ErrInsufficientBalance
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "invalid json body",
			body:       `not-json`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockTransactionService{createFn: tc.mockCreate}
			mux := newTxMux(svc)

			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(tc.body))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			if tc.wantAmount != 0 {
				var got models.Transaction
				assert.NoError(t, json.NewDecoder(w.Body).Decode(&got))
				assert.Equal(t, tc.wantAmount, got.Amount)
				assert.NotZero(t, got.ID)
			}
		})
	}
}
