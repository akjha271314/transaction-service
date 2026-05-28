package testutil

import (
	"database/sql"
	"testing"

	"transaction-service/internal/db"
)

// NewTestDB returns an initialised in-memory database that is closed
// automatically when the test finishes.
func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()
	database := db.Init()
	t.Cleanup(func() { database.Close() })
	return database
}
