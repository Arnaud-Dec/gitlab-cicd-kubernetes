package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"todoapp/internal/handlers"
	"todoapp/internal/version"
)

func TestVersion(t *testing.T) {
	original := version.Version
	version.Version = "2026-01-02T03:04"
	t.Cleanup(func() { version.Version = original })

	e := echo.New()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, handlers.Version(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"version":"2026-01-02T03:04"}`, rec.Body.String())
}
