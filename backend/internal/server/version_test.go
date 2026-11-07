package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"todoapp/internal/server"
	"todoapp/internal/version"
)

// The equivalent assertion for server.New lives in
// TestNew_SetsHostnameHeaderAndServesVersion (hostname_test.go) — see the
// comment there for why server.New can only be called once in this binary.
func TestNewHealthServer_ServesVersion(t *testing.T) {
	original := version.Version
	version.Version = "2026-01-02T03:04"
	t.Cleanup(func() { version.Version = original })

	e := server.NewHealthServer(nil)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"version":"2026-01-02T03:04"}`, rec.Body.String())
}
