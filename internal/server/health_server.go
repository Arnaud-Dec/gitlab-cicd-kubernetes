package server

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"todoapp/internal/handlers"
)

// NewHealthServer builds a separate Echo instance, listening on its own TCP
// port (HEALTH_PORT), that serves /healthz, /readyz, /version and /metrics
// for cluster-internal use only — Kubernetes liveness/readiness probes and
// Prometheus scraping. This port is never exposed outside the cluster, so it
// has no CORS setup; the frontend instead uses the copy of /healthz,
// /readyz and /version served on the main API port (see New in router.go).
func NewHealthServer(db *gorm.DB) *echo.Echo {
	e := echo.New()
	e.Use(hostnameMiddleware())

	healthHandler := handlers.NewHealthHandler(db)
	e.GET("/healthz", healthHandler.Live)
	e.GET("/readyz", healthHandler.Ready)
	e.GET("/version", handlers.Version)

	// Serves both the echo_* HTTP metrics registered by the API server's
	// echoprometheus middleware and the app's own todoapp_* metrics (see
	// internal/metrics) — both live in the default Prometheus registry.
	e.GET("/metrics", echoprometheus.NewHandler())

	return e
}
