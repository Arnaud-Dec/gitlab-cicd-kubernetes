package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"todoapp/internal/config"
	"todoapp/internal/models"
)

const (
	maxConnectAttempts = 10
	retryDelay         = 2 * time.Second
)

// Connect opens a connection to Postgres, retrying briefly to tolerate a
// database container that hasn't finished starting yet, then runs
// AutoMigrate for the app's models.
func Connect(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	var db *gorm.DB
	var err error

	for attempt := 1; attempt <= maxConnectAttempts; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
		if err == nil {
			break
		}
		log.Printf("database connection attempt %d/%d failed: %v", attempt, maxConnectAttempts, err)
		time.Sleep(retryDelay)
	}
	if err != nil {
		return nil, fmt.Errorf("could not connect to database after %d attempts: %w", maxConnectAttempts, err)
	}

	if err := db.AutoMigrate(&models.Todo{}); err != nil {
		return nil, fmt.Errorf("auto-migration failed: %w", err)
	}

	return db, nil
}
