# Backend — CI Guide

Go + Echo + GORM REST API for the Todo app. This file documents how to build, lint, and test this
project non-interactively — e.g. as CI pipeline steps.

## Prerequisites

- Go 1.27+ (`go.mod` pins `go 1.27.1`; 1.24+ is the actual minimum needed for the `tool` directive
  used to track `swag` and `golangci-lint` as build tools — see below — but this repo requires the
  newer version specifically)
- A Postgres 16 instance — only for the `test-postgres` and `run` commands below. `build`, `vet`,
  `lint`, and the default `test` need nothing external.

## Install dependencies

```bash
go mod download
```

This also makes the tracked build tools (`swag`, `golangci-lint`) resolvable via `go tool` — their
first invocation compiles them, which adds a few seconds the first time only.

## Build

```bash
# Plain build — Version falls back to the process's start time (see internal/version)
go build -o bin/api ./cmd/api

# Stamped with an explicit build timestamp instead — recommended for a real CI artifact
go build -ldflags "-X todoapp/internal/version.Version=$(date -u +%Y-%m-%dT%H:%M)" -o bin/api ./cmd/api
```

Produces a single binary at `bin/api` (gitignored) — the artifact to carry into a later
Docker-build or deploy stage.

## Vet and lint

```bash
go vet ./...
go tool golangci-lint run ./...
```

Config: `.golangci.yml`. Both exit non-zero on any finding — run them as separate CI steps so a
lint failure is reported distinctly from a vet failure.

## Test

```bash
# Default: in-memory SQLite, no external services — safe to run in any CI job as-is
go test ./...
```

To also run the exact same suite against real Postgres (e.g. as a second CI job with a Postgres
service container):

```bash
docker run -d --name todos-postgres -p 5432:5432 \
  -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=todos \
  postgres:16

# wait for it to accept connections, then:
TEST_DATABASE_DSN="host=localhost port=5432 user=postgres password=postgres dbname=todos sslmode=disable" \
  go test ./...
```

Both commands exit non-zero on any test failure. `-v` is optional, for more readable logs.

## Run (to feed a later E2E/integration stage)

```bash
# with a Postgres instance reachable (e.g. the container above)
go run ./cmd/api
# or, using the binary built above:
./bin/api
```

Listens on `HTTP_PORT` (default `8080`) and `HEALTH_PORT` (default `9090`). Poll
`GET :8080/readyz` (or `:9090/readyz`) until it returns `200` before pointing a later stage (the
`e2e/` or `cypress/` sub-projects) at it — it retries its own DB connection internally, but a CI
pipeline should still wait rather than assume it's immediately ready.

If that later stage is a **browser** hitting the API from a different origin (e.g. a served
frontend build, or the Cypress suite) rather than a plain HTTP client (e.g. the Go `e2e/` suite),
set `CORS_ENABLED=true` — see below. It's off by default, so a browser client gets blocked with no
extra step needed to get that default.

## Environment variables

| Var | Default | Notes |
|---|---|---|
| `DB_HOST` | `localhost` | |
| `DB_PORT` | `5432` | |
| `DB_USER` | `postgres` | |
| `DB_PASSWORD` | `postgres` | |
| `DB_NAME` | `todos` | |
| `DB_SSLMODE` | `disable` | |
| `HTTP_PORT` | `8080` | Main API port — the only one meant to be exposed outside the pipeline/cluster |
| `HEALTH_PORT` | `9090` | `/healthz`, `/readyz`, `/metrics` — keep this unexposed in a real deployment |
| `FRONTEND_ORIGIN` | `http://localhost:5173` | CORS allow-origin on `HTTP_PORT` — only used if `CORS_ENABLED` is true |
| `CORS_ENABLED` | `false` | Whether the API sends CORS headers at all — disabled means same-origin only, no `Access-Control-*` headers sent, and a cross-origin browser client is blocked by the browser itself even though the API is reachable |
| `TEST_DATABASE_DSN` | unset | Test-only — see "Test" above; when unset, tests use in-memory SQLite instead |

None of these need to be set for `build`, `vet`, `lint`, or the default `test` — they only matter
for `run` and `test-postgres`.

## Generate the OpenAPI spec (optional CI step)

```bash
go tool swag init -g cmd/api/main.go -o docs --outputTypes json
go run ./cmd/openapi
```

Writes `docs/openapi.json` and `docs/openapi.yaml` (gitignored) — a useful CI artifact to publish,
or to diff against a previous run to catch unintended API changes.
