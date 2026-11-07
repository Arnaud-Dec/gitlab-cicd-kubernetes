package server

// Internal test package (not server_test) so it can reach the unexported
// corsMiddleware helper directly — this avoids a second server.New call in
// this test binary, which would panic on duplicate Prometheus collector
// registration (see the comment on TestNew_SetsHostnameHeaderAndServesVersion
// in hostname_test.go). The wildcard CORS-disabled default is covered there
// instead.

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware_AllowsConfiguredOriginAndExposesHostnameHeader(t *testing.T) {
	e := echo.New()
	e.Use(corsMiddleware("http://localhost:5173", true))
	e.GET("/api/todos", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/todos", nil)
	req.Header.Set(echo.HeaderOrigin, "http://localhost:5173")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, "http://localhost:5173", rec.Header().Get(echo.HeaderAccessControlAllowOrigin))
	assert.Equal(t, HostnameHeader, rec.Header().Get(echo.HeaderAccessControlExposeHeaders))
}

func TestCORSMiddleware_AllowsEachOriginInACommaSeparatedList(t *testing.T) {
	e := echo.New()
	e.Use(corsMiddleware("http://localhost:5173, http://localhost:4000", true))
	e.GET("/api/todos", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	for _, origin := range []string{"http://localhost:5173", "http://localhost:4000"} {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/todos", nil)
		req.Header.Set(echo.HeaderOrigin, origin)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, origin, rec.Header().Get(echo.HeaderAccessControlAllowOrigin))
	}
}

func TestCORSMiddleware_RejectsOtherOriginsWhenEnabled(t *testing.T) {
	e := echo.New()
	e.Use(corsMiddleware("http://localhost:5173", true))
	e.GET("/api/todos", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/todos", nil)
	req.Header.Set(echo.HeaderOrigin, "http://evil.example.com")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Empty(t, rec.Header().Get(echo.HeaderAccessControlAllowOrigin))
}

func TestCORSMiddleware_AllowsAnyOriginWhenDisabled(t *testing.T) {
	e := echo.New()
	// frontendOrigin is ignored when corsEnabled is false — the middleware
	// falls back to a wildcard regardless of its value.
	e.Use(corsMiddleware("http://localhost:5173", false))
	e.GET("/api/todos", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/todos", nil)
	req.Header.Set(echo.HeaderOrigin, "http://anything.example.com")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, "*", rec.Header().Get(echo.HeaderAccessControlAllowOrigin))
	assert.Equal(t, HostnameHeader, rec.Header().Get(echo.HeaderAccessControlExposeHeaders))
}
