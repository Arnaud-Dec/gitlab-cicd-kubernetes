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

A **browser** hitting the API from a different origin (e.g. a served frontend build, or the
Cypress suite) works with no configuration needed — `CORS_ENABLED` defaults to `false`, which
allows any origin. Set it to `true` (and `FRONTEND_ORIGIN`) instead if the pipeline needs to lock
the API down to a specific origin — see below. None of this matters to a plain HTTP client (e.g.
the Go `e2e/` suite), which is never subject to CORS in the first place.

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
| `FRONTEND_ORIGIN` | `http://localhost:5173` | CORS allow-origin(s) on `HTTP_PORT` — only used if `CORS_ENABLED` is true; comma-separated for more than one |
| `CORS_ENABLED` | `false` | Whether the API restricts CORS to `FRONTEND_ORIGIN` (`true`) or allows any origin (`false`, the default) |
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

## Docker

The `Dockerfile` is multi-stage, and **every check above also exists as a stage**, so the exact same
command runs locally and in CI — no script is duplicated in `.gitlab-ci.yml`.

```bash
docker build -t backend .                                 # the API image (last stage, `api`)
docker build --target lint --output type=cacheonly .      # golangci-lint, no image produced
docker build --target test .                              # test suite on in-memory SQLite
docker build --target swagger .                           # Swagger UI serving the generated spec
```

| Stage | What it does |
|---|---|
| `base` | Go toolchain + dependencies, shared by every other stage |
| `lint` | `golangci-lint` (pinned prebuilt binary — building it via `go tool` costs minutes on the runner) |
| `test` | `go test` against in-memory SQLite |
| `openapi` | generates `docs/openapi.json` from the code annotations |
| `swagger` | `swaggerapi/swagger-ui` serving that spec |
| `api` | final Alpine image with the compiled binary — the deployed one |

`--output type=cacheonly` is worth knowing for the lint stage: the checks run, but Docker does not
assemble an image nobody will use.

Tests against a **real Postgres** need the database alongside, hence a dedicated compose file:

```bash
docker compose -f compose.test.yml run --rm test
docker compose -f compose.test.yml down -v --rmi local    # always clean up: the runner is shared
```

The whole application (API, frontend, Postgres, Swagger UI) runs from the `ci-cd` repository:
`docker compose up --build` in `../ci-cd`.

## CI pipeline

`.gitlab-ci.yml`, stages `lint` → `test` → `build` → `deploy`. Jobs use `docker:cli` and call the
Dockerfile stages above.

| Job | Stage | What it produces |
|---|---|---|
| `lint` | `lint` | nothing — fails the pipeline on a lint error |
| `test` | `test` | nothing — runs the suite against Postgres via `compose.test.yml` |
| `build-image` | `build` | `registry.granux.fr/decourt-demaria/tp/backend` |
| `build-swagger-image` | `build` | `registry.granux.fr/decourt-demaria/tp/backend/swagger` |
| `deploy` | `deploy` | triggers the CD pipeline of the `ci-cd` project |

Both image jobs share the same `.build-and-push` template and only differ by two variables
(`TARGET`, `IMAGE_REPO`). Every build is pushed as `:<short SHA>`, plus `:latest` on the default
branch and `:<git tag>` on a tag — so a deployment can always name the exact commit it runs.

Two details that cost time to rediscover:

- `docker build --provenance=false` is required: since Docker 29.7 the attestation can reach the
  registry before the image it refers to, which GitLab rejects with `blob unknown to registry`;
- the `test` job sets a unique `COMPOSE_PROJECT_NAME`, because the Docker daemon is shared with the
  other jobs running on the same runner.

The `deploy` job is a `trigger` job: it starts the pipeline of the `ci-cd` project with
`BACKEND_TAG: $CI_COMMIT_SHORT_SHA`, so the CD deploys the image this pipeline just built, and
`strategy: depend` makes this job red if the deployment fails.

## Deployment

Kubernetes manifests, the Helm chart and the CD are documented in the
[`ci-cd`](../ci-cd) directory. Two things differ from the
defaults in the table above once deployed there:

- `DB_USER` and `DB_PASSWORD` come from the secret generated by the Zalando Postgres operator, so
  no credential is ever written in a manifest;
- `DB_SSLMODE` must be `require` — the operator's Postgres image refuses unencrypted connections.
