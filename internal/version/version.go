// Package version holds the build-time version stamp for the backend
// binary. It's a build timestamp (e.g. "2026-09-06T21:50"), set via
// `-ldflags "-X todoapp/internal/version.Version=..."` — see the README's
// "Build with a version stamp" section.
package version

import "time"

// Version defaults to "dev": the -X linker flag only works on a variable
// with a plain string-literal initializer, so it can't be given a computed
// default directly. Call Init at startup to fall back to the current time
// whenever ldflags weren't used (e.g. a plain `go run`/`go build`) — real
// CI/production builds set this via ldflags instead, which Init leaves
// untouched.
var Version = "dev"

// Init falls back Version to the current UTC time (to the minute) if it's
// still at its zero-ldflags default. Safe to call unconditionally at
// startup — a real -X-stamped Version is left exactly as given.
func Init() {
	if Version == "dev" {
		Version = time.Now().UTC().Format("2006-01-02T15:04")
	}
}
