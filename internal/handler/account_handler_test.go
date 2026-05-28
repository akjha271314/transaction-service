package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"transaction-service/internal/models"
	"transaction-service/internal/service"
)

type mockAccountService struct {
	createFn func(documentNumber string, balance float64) (*models.Account, error)
	getFn    func(id int64) (*models.Account, error)
}

func (m *mockAccountService) CreateAccount(documentNumber string, balance float64) (*models.Account, error) {
	return m.createFn(documentNumber, balance)
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

func TestAccountHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mockCreate func(string, float64) (*models.Account, error)
		wantStatus int
		wantBody   *models.Account
	}{
		{
			name: "success with balance",
			body: `{"document_number":"12345678900","balance":500}`,
			mockCreate: func(doc string, bal float64) (*models.Account, error) {
				return &models.Account{ID: 1, DocumentNumber: doc, Balance: bal}, nil
			},
			wantStatus: http.StatusCreated,
			wantBody:   &models.Account{ID: 1, DocumentNumber: "12345678900", Balance: 500},
		},
		{
			name: "success with zero balance",
			body: `{"document_number":"12345678900"}`,
			mockCreate: func(doc string, bal float64) (*models.Account, error) {
				return &models.Account{ID: 2, DocumentNumber: doc, Balance: 0}, nil
			},
			wantStatus: http.StatusCreated,
			wantBody:   &models.Account{ID: 2, DocumentNumber: "12345678900", Balance: 0},
		},
		{
			name: "duplicate document number returns 409",
			body: `{"document_number":"12345678900","balance":100}`,
			mockCreate: func(doc string, bal float64) (*models.Account, error) {
				return nil, service.ErrDuplicateAccount
			},
			wantStatus: http.StatusConflict,
		},
		{
			name:       "missing document number",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "negative balance",
			body:       `{"document_number":"123","balance":-100}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       `not-json`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockAccountService{createFn: tc.mockCreate}
			mux := newAccountMux(svc)

			req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString(tc.body))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			if tc.wantBody != nil {
				var got models.Account
				require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
				assert.Equal(t, tc.wantBody.ID, got.ID)
				assert.Equal(t, tc.wantBody.DocumentNumber, got.DocumentNumber)
				assert.Equal(t, tc.wantBody.Balance, got.Balance)
			}
		})
	}
}

func TestAccountHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		mockGet    func(int64) (*models.Account, error)
		wantStatus int
		wantID     int64
	}{
		{
			name: "existing account returns full data",
			url:  "/accounts/1",
			mockGet: func(id int64) (*models.Account, error) {
				return &models.Account{ID: id, DocumentNumber: "12345678900", Balance: 500}, nil
			},
			wantStatus: http.StatusOK,
			wantID:     1,
		},
		{
			name: "account not found returns 404",
			url:  "/accounts/99",
			mockGet: func(id int64) (*models.Account, error) {
				return nil, service.ErrAccountNotFound
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "non-numeric id returns 400",
			url:        "/accounts/abc",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockAccountService{getFn: tc.mockGet}
			mux := newAccountMux(svc)

			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			if tc.wantID != 0 {
				var got models.Account
				require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
				assert.Equal(t, tc.wantID, got.ID)
				assert.Equal(t, "12345678900", got.DocumentNumber)
				assert.Equal(t, float64(500), got.Balance)
			}
		})
	}
}
