package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"transaction-service/internal/models"
	"transaction-service/internal/repository"
)

func TestAccountService_CreateAccount(t *testing.T) {
	tests := []struct {
		name           string
		documentNumber string
		balance        float64
		mockErr        error
		wantErr        error
		wantBalance    float64
	}{
		{
			name:           "success with positive balance",
			documentNumber: "12345678900",
			balance:        500.0,
			wantBalance:    500.0,
		},
		{
			name:           "success with zero balance",
			documentNumber: "99999999999",
			balance:        0.0,
			wantBalance:    0.0,
		},
		{
			name:           "duplicate document number",
			documentNumber: "12345678900",
			balance:        100.0,
			mockErr:        repository.ErrDuplicateAccount,
			wantErr:        ErrDuplicateAccount,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockAccountRepo{
				createFn: func(documentNumber string, balance float64) (*models.Account, error) {
					if tc.mockErr != nil {
						return nil, tc.mockErr
					}
					return &models.Account{ID: 1, DocumentNumber: documentNumber, Balance: balance}, nil
				},
			}
			svc := NewAccountService(repo)

			acc, err := svc.CreateAccount(tc.documentNumber, tc.balance)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, acc)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.documentNumber, acc.DocumentNumber)
				assert.Equal(t, tc.wantBalance, acc.Balance)
				assert.NotZero(t, acc.ID)
			}
		})
	}
}

func TestAccountService_GetAccount(t *testing.T) {
	tests := []struct {
		name      string
		accountID int64
		mockFn    func(id int64) (*models.Account, error)
		wantErr   error
		wantID    int64
	}{
		{
			name:      "existing account returns correct data",
			accountID: 1,
			mockFn: func(id int64) (*models.Account, error) {
				return &models.Account{ID: id, DocumentNumber: "12345678900", Balance: 500}, nil
			},
			wantID: 1,
		},
		{
			name:      "non-existent account returns ErrAccountNotFound",
			accountID: 99,
			mockFn: func(id int64) (*models.Account, error) {
				return nil, repository.ErrNotFound
			},
			wantErr: ErrAccountNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockAccountRepo{findByIDFn: tc.mockFn}
			svc := NewAccountService(repo)

			acc, err := svc.GetAccount(tc.accountID)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, acc)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantID, acc.ID)
				assert.NotEmpty(t, acc.DocumentNumber)
			}
		})
	}
}
