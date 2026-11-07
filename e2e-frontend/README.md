# Cypress (Frontend E2E) — CI Guide

Browser-driven tests against a **running** frontend (which itself needs a running backend and
Postgres), using [Cypress](https://www.cypress.io) — a separate npm project, not part of the
frontend's own build. Unlike the frontend's own Vitest tests (which mount components in isolation
with the API mocked), these drive a real browser against the actual running app. This file
documents how to run it headlessly in CI.

## Prerequisites

- Node.js 20+
- The **full stack running**: Postgres, the backend (`backend/`), and the frontend (`frontend/`,
  dev server or a served production build) — this suite drives a real browser against them; it
  starts none of them itself. In a CI pipeline these are prior steps or jobs.
- The backend must be started with **`CORS_ENABLED=true`**. It's disabled by default, and since
  this suite is a real browser hitting the frontend's origin (not the backend's), every API call
  would otherwise be blocked by the browser itself even though the backend is up and reachable.
- Cypress bundles its own browser (Electron) — no separate browser install needed for the default
  setup. On a minimal CI image (e.g. a bare Docker container) you may still need a few system
  libraries Cypress depends on; see
  [Cypress's own CI documentation](https://docs.cypress.io/app/continuous-integration/overview) if
  `cypress run` fails to *launch* rather than fails a test.

## Install dependencies

```bash
npm ci
```

## Run (headless — the CI-safe mode)

```bash
npx cypress run
```

`npm test` does the same but first blocks on `wait-on` until the frontend responds — convenient
locally, redundant if your CI pipeline already waits for the frontend service to be healthy before
starting this step.

Defaults to targeting `http://localhost:5173`. Point it at wherever the frontend from a prior step
is actually reachable with `CYPRESS_BASE_URL` (Cypress reads any `CYPRESS_`-prefixed env var
automatically):

```bash
CYPRESS_BASE_URL=http://frontend:5173 npx cypress run
```

## Environment variables

| Var | Default | Notes |
|---|---|---|
| `CYPRESS_BASE_URL` | `http://localhost:5173` | Base URL of the running frontend to test against |

Exits non-zero if any spec fails.

## Artifacts on failure

Screenshots (gitignored, under `cypress/screenshots/`) are captured automatically on a failing
test — have your CI pipeline upload them as build artifacts so a failure is debuggable from the CI
UI, not just its text log. Video recording is off by default in this project's Cypress version, so
don't expect files under `cypress/videos/` unless you explicitly enable it in `cypress.config.js`.

## Test data

Every spec that creates a todo deletes it again in `afterEach` (via the idempotent
`cy.deleteTodoByTitle` custom command) — safe to run repeatedly against the same persistent
backend/database without manual cleanup between CI runs.
