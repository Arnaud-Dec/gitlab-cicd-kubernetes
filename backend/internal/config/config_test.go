package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"todoapp/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	t.Chdir(t.TempDir()) // no .env file here, and no relevant env vars set

	cfg := config.Load()

	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "postgres", cfg.DBPassword)
	assert.Equal(t, "todos", cfg.DBName)
	assert.Equal(t, "disable", cfg.DBSSLMode)
	assert.Equal(t, "8080", cfg.HTTPPort)
	assert.Equal(t, "9090", cfg.HealthPort)
	assert.Equal(t, "http://localhost:5173", cfg.FrontendOrigin)
	assert.False(t, cfg.CORSEnabled)
}

func TestLoad_CORSEnabledCanBeTurnedOn(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("CORS_ENABLED", "true")

	cfg := config.Load()

	assert.True(t, cfg.CORSEnabled)
}

func TestLoad_EnvVarOverridesDefault(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("HTTP_PORT", "3000")

	cfg := config.Load()

	assert.Equal(t, "3000", cfg.HTTPPort)
}

func TestLoad_DotEnvFileOverridesDefault(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	envFile := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envFile, []byte("DB_NAME=from_dotenv\n"), 0o600))

	cfg := config.Load()

	assert.Equal(t, "from_dotenv", cfg.DBName)
}

func TestLoad_EnvVarOverridesDotEnvFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	envFile := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envFile, []byte("DB_NAME=from_dotenv\n"), 0o600))
	t.Setenv("DB_NAME", "from_os_env")

	cfg := config.Load()

	assert.Equal(t, "from_os_env", cfg.DBName)
}
