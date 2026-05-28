package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"transaction-service/internal/testutil"
)

func TestTransactionRepository_CreateTx(t *testing.T) {
	tests := []struct {
		name            string
		operationTypeID int64
		amount          float64
	}{
		{"purchase (negative amount)", 1, -50.0},
		{"credit voucher (positive amount)", 4, 60.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := testutil.NewTestDB(t)
			accountRepo := NewAccountRepository(db)
			txRepo := NewTransactionRepository(db)

			acc, err := accountRepo.Create("12345678900", 500.0)
			require.NoError(t, err)

			sqlTx, err := db.Begin()
			require.NoError(t, err)
			defer sqlTx.Rollback()

			tx, err := txRepo.CreateTx(sqlTx, acc.ID, tc.operationTypeID, tc.amount)
			require.NoError(t, err)
			sqlTx.Commit()

			assert.NotZero(t, tx.ID)
			assert.Equal(t, acc.ID, tx.AccountID)
			assert.Equal(t, tc.operationTypeID, tx.OperationTypeID)
			assert.Equal(t, tc.amount, tx.Amount)
			assert.False(t, tx.EventDate.IsZero())
		})
	}
}

func TestTransactionRepository_CreateTx_InvalidAccount(t *testing.T) {
	db := testutil.NewTestDB(t)
	txRepo := NewTransactionRepository(db)

	sqlTx, err := db.Begin()
	require.NoError(t, err)
	defer sqlTx.Rollback()

	_, err = txRepo.CreateTx(sqlTx, 999, 1, -50.0)
	assert.Error(t, err)
}

func TestTransactionRepository_FindOperationType(t *testing.T) {
	tests := []struct {
		id          int64
		description string
		isCredit    bool
	}{
		{1, "Normal Purchase", false},
		{2, "Purchase with installments", false},
		{3, "Withdrawal", false},
		{4, "Credit Voucher", true},
	}

	db := testutil.NewTestDB(t)
	txRepo := NewTransactionRepository(db)

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			op, err := txRepo.FindOperationType(tc.id)

			require.NoError(t, err)
			assert.Equal(t, tc.id, op.ID)
			assert.Equal(t, tc.description, op.Description)
			assert.Equal(t, tc.isCredit, op.IsCredit)
		})
	}
}

func TestTransactionRepository_FindOperationType_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	txRepo := NewTransactionRepository(db)

	_, err := txRepo.FindOperationType(99)
	assert.ErrorIs(t, err, ErrNotFound)
}
