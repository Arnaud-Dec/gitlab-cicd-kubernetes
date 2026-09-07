package server

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

// HostnameHeader is the response header carrying the server's hostname — in
// Kubernetes this is the pod name by default, letting a client see which
// replica served a given request (e.g. to watch a Service load-balance
// across pods, or confirm a rolling update reached a given instance).
const HostnameHeader = "X-Server-Hostname"

func hostnameMiddleware() echo.MiddlewareFunc {
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("could not determine hostname: %v", err)
		hostname = "unknown"
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(HostnameHeader, hostname)
			return next(c)
		}
	}
}
