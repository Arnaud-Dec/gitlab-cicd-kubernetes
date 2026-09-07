package server

import (
	"log"

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
// corsEnabled toggles CORS support (CORS_ENABLED, default false — see
// config.Config). It's off by default so the API is only reachable
// same-origin unless a deployment deliberately opts in; local dev, where the
// frontend (:5173) and backend (:8080) are different origins, needs it on.
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
	if corsEnabled {
		e.Use(corsMiddleware(frontendOrigin))
	}
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

func corsMiddleware(frontendOrigin string) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontendOrigin},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
		// Custom response headers are invisible to cross-origin JS unless
		// explicitly exposed — without this, HostnameHeader would show up
		// fine in curl/devtools' network tab but never reach axios.
		ExposeHeaders: []string{HostnameHeader},
	})
}
