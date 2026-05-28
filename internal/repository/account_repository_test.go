package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"transaction-service/internal/testutil"
)

func TestAccountRepository_Create(t *testing.T) {
	tests := []struct {
		name           string
		documentNumber string
		balance        float64
	}{
		{"with balance", "12345678900", 500.0},
		{"with zero balance", "99999999999", 0.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := testutil.NewTestDB(t)
			repo := NewAccountRepository(db)

			acc, err := repo.Create(tc.documentNumber, tc.balance)

			require.NoError(t, err)
			assert.NotZero(t, acc.ID)
			assert.Equal(t, tc.documentNumber, acc.DocumentNumber)
			assert.Equal(t, tc.balance, acc.Balance)
		})
	}
}

func TestAccountRepository_Create_DuplicateDocumentNumber(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	_, err := repo.Create("12345678900", 0)
	require.NoError(t, err)

	_, err = repo.Create("12345678900", 0)
	assert.ErrorIs(t, err, ErrDuplicateAccount)
}

func TestAccountRepository_FindByID(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	created, err := repo.Create("12345678900", 1000.0)
	require.NoError(t, err)

	found, err := repo.FindByID(created.ID)

	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, created.DocumentNumber, found.DocumentNumber)
	assert.Equal(t, created.Balance, found.Balance)
}

func TestAccountRepository_FindByID_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewAccountRepository(db)

	_, err := repo.FindByID(999)

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestAccountRepository_UpdateBalanceTx(t *testing.T) {
	tests := []struct {
		name            string
		initialBalance  float64
		delta           float64
		wantBalance     float64
		wantErr         error
	}{
		{"debit within balance", 100.0, -40.0, 60.0, nil},
		{"debit exact balance", 100.0, -100.0, 0.0, nil},
		{"credit from zero", 0.0, 100.0, 100.0, nil},
		{"credit adds to existing", 50.0, 25.0, 75.0, nil},
		{"debit exceeds balance", 30.0, -50.0, 0, ErrInsufficientBalance},
		{"debit on zero balance", 0.0, -10.0, 0, ErrInsufficientBalance},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := testutil.NewTestDB(t)
			repo := NewAccountRepository(db)

			acc, err := repo.Create("12345678900", tc.initialBalance)
			require.NoError(t, err)

			sqlTx, err := db.Begin()
			require.NoError(t, err)
			defer sqlTx.Rollback()

			err = repo.UpdateBalanceTx(sqlTx, acc.ID, tc.delta)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				sqlTx.Commit()
				updated, _ := repo.FindByID(acc.ID)
				assert.Equal(t, tc.wantBalance, updated.Balance)
			}
		})
	}
}
