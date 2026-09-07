package main

import (
	"errors"
	"log"
	"net/http"

	"todoapp/internal/config"
	"todoapp/internal/database"
	"todoapp/internal/metrics"
	"todoapp/internal/server"
	"todoapp/internal/version"
)

// @title			Todo App API
// @version		1.0
// @description	REST API for a simple Todo/task manager, used as a teaching reference app for a DevOps course.
// @description	/healthz and /readyz below are also mirrored on a separate, unexposed port (HEALTH_PORT, default 9090) reserved for cluster-internal use (Kubernetes probes, metrics).
//
// @BasePath	/
func main() {
	version.Init()

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := metrics.SeedFromDB(db); err != nil {
		log.Fatalf("failed to seed metrics from database: %v", err)
	}

	healthServer := server.NewHealthServer(db)
	go func() {
		log.Printf("health server listening on :%s", cfg.HealthPort)
		if err := healthServer.Start(":" + cfg.HealthPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("health server stopped: %v", err)
		}
	}()

	e := server.New(db, cfg.FrontendOrigin, cfg.CORSEnabled)

	log.Printf("listening on :%s", cfg.HTTPPort)
	if err := e.Start(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
