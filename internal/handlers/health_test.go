package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"todoapp/internal/handlers"
)

func TestHealthHandler_Live(t *testing.T) {
	db := setupTestDB(t)
	h := handlers.NewHealthHandler(db)

	e := echo.New()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Live(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHealthHandler_Ready_DBUp(t *testing.T) {
	db := setupTestDB(t)
	h := handlers.NewHealthHandler(db)

	e := echo.New()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Ready(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHealthHandler_Ready_DBDown(t *testing.T) {
	db := setupTestDB(t)
	h := handlers.NewHealthHandler(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	e := echo.New()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Ready(c))
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}
