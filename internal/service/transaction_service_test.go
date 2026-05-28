package service

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"transaction-service/internal/models"
	"transaction-service/internal/repository"
)

func TestTransactionService_CreateTransaction_Success(t *testing.T) {
	txRepo := &mockTransactionRepo{
		findOperationTypeFn: func(id int64) (*models.OperationType, error) {
			return &models.OperationType{ID: id, IsCredit: false}, nil
		},
		createTxFn: func(tx *sql.Tx, accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
			return &models.Transaction{
				ID:              1,
				AccountID:       accountID,
				OperationTypeID: operationTypeID,
				Amount:          amount,
				EventDate:       time.Now(),
			}, nil
		},
	}
	svc := newSvc(txRepo, makeAccountRepo(&models.Account{ID: 1, Balance: 500}, nil))

	tx, err := svc.CreateTransaction(1, 1, 50.0)

	require.NoError(t, err)
	assert.NotZero(t, tx.ID)
	assert.Equal(t, int64(1), tx.AccountID)
	assert.Equal(t, int64(1), tx.OperationTypeID)
	assert.Equal(t, -50.0, tx.Amount)
	assert.False(t, tx.EventDate.IsZero())
}

func TestTransactionService_SignLogic(t *testing.T) {
	tests := []struct {
		name        string
		opTypeID    int64
		isCredit    bool
		inputAmount float64
		wantAmount  float64
	}{
		{"normal purchase stores negative", 1, false, 50.0, -50.0},
		{"withdrawal stores negative", 3, false, 20.0, -20.0},
		{"credit voucher stores positive", 4, true, 60.0, 60.0},
		{"negative input for debit is normalised", 1, false, -50.0, -50.0},
		{"negative input for credit is normalised", 4, true, -60.0, 60.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var stored float64
			txRepo := &mockTransactionRepo{
				findOperationTypeFn: func(id int64) (*models.OperationType, error) {
					return &models.OperationType{ID: id, IsCredit: tc.isCredit}, nil
				},
				createTxFn: func(tx *sql.Tx, accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
					stored = amount
					return &models.Transaction{Amount: amount}, nil
				},
			}
			svc := newSvc(txRepo, makeAccountRepo(&models.Account{ID: 1, Balance: 1000}, nil))

			result, err := svc.CreateTransaction(1, tc.opTypeID, tc.inputAmount)

			require.NoError(t, err)
			assert.Equal(t, tc.wantAmount, stored)
			assert.Equal(t, tc.wantAmount, result.Amount)
		})
	}
}

func TestTransactionService_CreateTransaction_Errors(t *testing.T) {
	tests := []struct {
		name        string
		accountRepo *mockAccountRepo
		txRepo      *mockTransactionRepo
		wantErr     error
	}{
		{
			name:        "invalid account returns ErrInvalidAccount",
			accountRepo: makeAccountRepo(nil, repository.ErrNotFound),
			txRepo:      makeTxRepo(nil),
			wantErr:     ErrInvalidAccount,
		},
		{
			name:        "invalid operation type returns ErrInvalidOperationType",
			accountRepo: makeAccountRepo(&models.Account{ID: 1}, nil),
			txRepo: &mockTransactionRepo{
				findOperationTypeFn: func(id int64) (*models.OperationType, error) {
					return nil, repository.ErrNotFound
				},
			},
			wantErr: ErrInvalidOperationType,
		},
		{
			name:        "insufficient balance returns ErrInsufficientBalance",
			accountRepo: insufficientBalanceRepo(),
			txRepo: &mockTransactionRepo{
				findOperationTypeFn: func(id int64) (*models.OperationType, error) {
					return &models.OperationType{ID: id, IsCredit: false}, nil
				},
				createTxFn: func(tx *sql.Tx, a, b int64, c float64) (*models.Transaction, error) {
					return nil, nil
				},
			},
			wantErr: ErrInsufficientBalance,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := newSvc(tc.txRepo, tc.accountRepo)

			_, err := svc.CreateTransaction(1, 1, 50.0)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
