package server_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"todoapp/internal/server"
	"todoapp/internal/version"
)

// Live (/healthz) never touches the database, so a nil *gorm.DB is safe here
// — these tests only care about the hostname header, set by middleware
// before any handler runs.

// echoprometheus.NewMiddleware (used by server.New) registers its collectors
// into the default Prometheus registry, which panics on a second
// registration — so server.New must only be called once across this whole
// test binary. This test therefore also covers /version and the CORS-off
// default (see version_test.go for the equivalent on NewHealthServer, which
// doesn't have this constraint, and cors_internal_test.go for coverage of
// the CORS-on case, which tests the extracted middleware directly instead of
// going through another server.New call).
func TestNew_SetsHostnameHeaderAndServesVersion(t *testing.T) {
	wantHostname, err := os.Hostname()
	require.NoError(t, err)

	originalVersion := version.Version
	version.Version = "2026-01-02T03:04"
	t.Cleanup(func() { version.Version = originalVersion })

	// corsEnabled: false — the actual default (config.Config.CORSEnabled).
	e := server.New(nil, "http://localhost:5173", false)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, wantHostname, rec.Header().Get(server.HostnameHeader))

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/version", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"version":"2026-01-02T03:04"}`, rec.Body.String())

	// With CORS disabled, a cross-origin request gets no CORS headers at all
	// — the browser will block the frontend from reading the response. Using
	// /healthz here (rather than /api/todos) since it doesn't touch the nil
	// *gorm.DB this test passes to server.New.
	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/healthz", nil)
	req.Header.Set(echo.HeaderOrigin, "http://localhost:5173")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Empty(t, rec.Header().Get(echo.HeaderAccessControlAllowOrigin))
}

func TestNewHealthServer_SetsHostnameHeader(t *testing.T) {
	wantHostname, err := os.Hostname()
	require.NoError(t, err)

	e := server.NewHealthServer(nil)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, wantHostname, rec.Header().Get(server.HostnameHeader))
}
