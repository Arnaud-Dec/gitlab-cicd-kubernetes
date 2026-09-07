# E2E (Go/Ginkgo) — CI Guide

Black-box HTTP tests against a **running** backend, using
[Ginkgo](https://onsi.github.io/ginkgo/)/[Gomega](https://onsi.github.io/gomega/) — a separate Go
module, not a package inside the backend itself. Unlike the backend's own unit tests (which call
handlers directly, in-process), these only ever talk to it over real HTTP, the same way a real
client would. This file documents how to run it non-interactively in CI.

## Prerequisites

- Go 1.27+ (`go.mod` pins `go 1.27.1`)
- A **running backend**, itself backed by Postgres. This suite makes real HTTP calls — it does not
  start the backend itself. In a CI pipeline, starting `backend/` (and its Postgres service) is a
  prior step or job, not something this sub-project owns.

## Install dependencies

```bash
go mod download
```

## Run

```bash
go test ./...
# or, for Ginkgo's own more detailed reporter:
go tool ginkgo ./...
```

Defaults to targeting `http://localhost:8080`. Point it at wherever the backend from a prior CI
step is actually reachable with `API_BASE_URL`:

```bash
API_BASE_URL=http://backend:8080 go test ./...
```

## Environment variables

| Var | Default | Notes |
|---|---|---|
| `API_BASE_URL` | `http://localhost:8080` | Base URL of the running backend to test against |

## Startup behavior and failure mode

The suite's `BeforeSuite` polls `GET $API_BASE_URL/readyz` for up to 30s before running any spec.
This tolerates the backend still finishing its own startup/DB-connection retry, but does **not**
replace a CI-level "wait for the service" step if the backend needs longer than that to come up
(e.g. a cold container pull) — that timeout will simply fail the whole suite instead.

Exits non-zero if any spec fails, or if the backend never became ready — in the latter case with a
single clear message naming the URL it tried, rather than 13 separate confusing
connection-refused errors.

## Test data

Every spec that creates a todo deletes it again (`AfterEach`/`DeferCleanup`) — safe to run
repeatedly against the same persistent backend/database without manual cleanup between CI runs.
