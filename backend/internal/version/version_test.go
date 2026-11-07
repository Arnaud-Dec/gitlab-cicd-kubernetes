package version_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"todoapp/internal/version"
)

func TestInit_FallsBackToCurrentTimeWhenUnset(t *testing.T) {
	original := version.Version
	version.Version = "dev"
	t.Cleanup(func() { version.Version = original })

	before := time.Now().UTC()
	version.Init()
	after := time.Now().UTC()

	stamped, err := time.Parse("2006-01-02T15:04", version.Version)
	assert.NoError(t, err)
	assert.WithinRange(t, stamped, before.Truncate(time.Minute), after.Truncate(time.Minute).Add(time.Minute))
}

func TestInit_LeavesAnLdflagsValueUntouched(t *testing.T) {
	original := version.Version
	version.Version = "2026-01-02T03:04"
	t.Cleanup(func() { version.Version = original })

	version.Init()

	assert.Equal(t, "2026-01-02T03:04", version.Version)
}
