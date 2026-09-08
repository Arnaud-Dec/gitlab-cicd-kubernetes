package server

import (
	"log"
	"strings"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"

	"todoapp/internal/handlers"
)

// New builds the Echo instance serving the Todo API, with middleware and
// routes wired up. /healthz, /readyz and /version are also served here (in
// addition to NewHealthServer's copy on HEALTH_PORT) because this is the
// only port exposed to the frontend — HEALTH_PORT stays unexposed, reserved
// for cluster-internal use (Kubernetes probes, metrics).
//
// CORS headers are always sent. corsEnabled (CORS_ENABLED, default false —
// see config.Config) picks which policy: false (the default) allows any
// origin (Access-Control-Allow-Origin: *), a permissive default so local dev
// and ad hoc testing work regardless of what port the frontend happens to
// run on; true restricts it to frontendOrigin instead, for deployments that
// want to lock the API down to a known set of origins.
//
// frontendOrigin accepts a comma-separated list (e.g. when the frontend dev
// server needs to run on more than one port at once), see splitOrigins. It's
// only consulted when corsEnabled is true.
func New(db *gorm.DB, frontendOrigin string, corsEnabled bool) *echo.Echo {
	e := echo.New()
	e.Validator = NewValidator()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:  true,
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,
		LogError:   true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				log.Printf("%s %s status=%d latency=%s error=%v", v.Method, v.URI, v.Status, v.Latency, v.Error)
			} else {
				log.Printf("%s %s status=%d latency=%s", v.Method, v.URI, v.Status, v.Latency)
			}
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(corsMiddleware(frontendOrigin, corsEnabled))
	e.Use(hostnameMiddleware())
	// Collects request count/duration/size metrics (echo_requests_total,
	// echo_request_duration_seconds, ...) into the default Prometheus
	// registry. They're scraped from the admin server's /metrics — see
	// NewHealthServer — not from this port.
	e.Use(echoprometheus.NewMiddleware(""))

	todoHandler := handlers.NewTodoHandler(db)
	healthHandler := handlers.NewHealthHandler(db)

	e.GET("/healthz", healthHandler.Live)
	e.GET("/readyz", healthHandler.Ready)
	e.GET("/version", handlers.Version)

	api := e.Group("/api")
	api.GET("/todos", todoHandler.List)
	api.POST("/todos", todoHandler.Create)
	api.GET("/todos/:id", todoHandler.Get)
	api.PUT("/todos/:id", todoHandler.Update)
	api.PATCH("/todos/:id/toggle", todoHandler.Toggle)
	api.DELETE("/todos/:id", todoHandler.Delete)

	return e
}

func corsMiddleware(frontendOrigin string, corsEnabled bool) echo.MiddlewareFunc {
	origins := []string{"*"}
	if corsEnabled {
		origins = splitOrigins(frontendOrigin)
	}
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
		// Custom response headers are invisible to cross-origin JS unless
		// explicitly exposed — without this, HostnameHeader would show up
		// fine in curl/devtools' network tab but never reach axios.
		ExposeHeaders: []string{HostnameHeader},
	})
}

// splitOrigins turns a comma-separated FRONTEND_ORIGIN value into the list
// AllowOrigins expects, trimming whitespace and dropping empty entries. A
// single origin (the common case) just comes back as a one-element slice.
func splitOrigins(frontendOrigin string) []string {
	parts := strings.Split(frontendOrigin, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			origins = append(origins, p)
		}
	}
	return origins
}
