package config

import "github.com/spf13/viper"

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	HTTPPort       string
	HealthPort     string
	FrontendOrigin string
	CORSEnabled    bool
}

func Load() Config {
	v := viper.New()

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "todos")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("HEALTH_PORT", "9090")
	v.SetDefault("FRONTEND_ORIGIN", "http://localhost:5173")
	v.SetDefault("CORS_ENABLED", false)

	// Optional .env file, same convenience godotenv used to provide —
	// ignored if absent. Real environment variables still take precedence
	// over it (and over the defaults above) via AutomaticEnv below.
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	_ = v.ReadInConfig()

	v.AutomaticEnv()

	return Config{
		DBHost:         v.GetString("DB_HOST"),
		DBPort:         v.GetString("DB_PORT"),
		DBUser:         v.GetString("DB_USER"),
		DBPassword:     v.GetString("DB_PASSWORD"),
		DBName:         v.GetString("DB_NAME"),
		DBSSLMode:      v.GetString("DB_SSLMODE"),
		HTTPPort:       v.GetString("HTTP_PORT"),
		HealthPort:     v.GetString("HEALTH_PORT"),
		FrontendOrigin: v.GetString("FRONTEND_ORIGIN"),
		CORSEnabled:    v.GetBool("CORS_ENABLED"),
	}
}
