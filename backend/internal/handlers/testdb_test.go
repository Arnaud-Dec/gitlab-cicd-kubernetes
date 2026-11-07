package handlers_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"todoapp/internal/models"
)

func itoa(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

// setupTestDB returns a fresh database for a single test. By default it uses
// an in-memory SQLite database, so `go test ./...` needs no external
// services. If TEST_DATABASE_DSN is set, it connects to that Postgres
// instance instead and clears the todos table, so the exact same test suite
// can also be run against real Postgres (e.g. from a CI service container).
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_DSN")

	gormConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}

	var db *gorm.DB
	var err error
	if dsn == "" {
		db, err = gorm.Open(sqlite.Open(":memory:"), gormConfig)
	} else {
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	}
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&models.Todo{}))

	if dsn != "" {
		require.NoError(t, db.Exec("TRUNCATE TABLE todos RESTART IDENTITY CASCADE").Error)
	}

	return db
}
