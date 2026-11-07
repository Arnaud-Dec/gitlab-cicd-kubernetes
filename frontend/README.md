# Frontend — CI Guide

Vue 3 + Vite SPA for the Todo app. This file documents how to install, test, and build this
project non-interactively — e.g. as CI pipeline steps.

## Prerequisites

- Node.js 20+

## Install dependencies

```bash
npm ci
```

Use `npm ci`, not `npm install`, in CI: it's deterministic (installs exactly what
`package-lock.json` specifies) and fails loudly if the lockfile is out of sync, rather than
silently updating it.

## Test

```bash
npm test
```

Runs the Vitest suite once and exits (`vitest run`, not watch mode — already CI-safe). **No
backend, database, or network access needed**: every API call is mocked at the module level (see
`src/__tests__/`). Exits non-zero on any failure.

## Build

```bash
npm run build
```

Produces the production bundle in `dist/` (gitignored) — the artifact to carry into a later
Docker-build or deploy stage. `dist/config.json` ships with localhost defaults and is meant to be
replaced post-build for a real deployment — that's a later exercise, not something this build step
does.

No environment variables are needed for either command: the frontend's backend address is resolved
at runtime from `config.json`, not baked in at build time, and its own version stamp (shown in the
status bar) is computed automatically from the current date — no input needed.

## Docker

The `Dockerfile` is multi-stage, and the test suite is one of its stages, so the same command runs
locally and in CI:

```bash
docker build -t frontend .          # production image: nginx serving dist/
docker build --target test .        # runs `npm test`; fails the build on a test failure
```

| Stage | What it does |
|---|---|
| `base` | `npm ci --ignore-scripts` + sources |
| `test` | `npm test` |
| `builder` | `npm run build` → `dist/` |
| final | `nginx:stable-alpine` serving `dist/` on **port 80** |

`npm ci --ignore-scripts` is deliberate: no dependency of this project needs a postinstall script,
and skipping them removes the attack vector used by the March 2026 axios compromise.

Note the port: the container listens on **80**, not on Vite's 5173. The `docker-compose.yaml` of the
`ci-cd` repository maps it to 5173 on the host, which is where the 5173 you may remember comes from.

The whole application (frontend, API, Postgres, Swagger UI) runs from the `ci-cd` repository:
`docker compose up --build` in `../ci-cd`.

## Backend address at runtime

`dist/config.json` is read by the browser at startup, not baked into the bundle — so **one image
works in every environment**. In Kubernetes it is replaced by a ConfigMap mounted over the file:

```yaml
volumeMounts:
  - name: config
    mountPath: /usr/share/nginx/html/config.json
    subPath: config.json
```

A file mounted with `subPath` is **not** refreshed when the ConfigMap changes, so after editing it
run `kubectl rollout restart deploy/<frontend>`. See the `ci-cd` repository for the full manifest.

## CI pipeline

`.gitlab-ci.yml`, stages `test` → `build` → `deploy`, all on `docker:cli`.

| Job | Stage | What it produces |
|---|---|---|
| `test` | `test` | nothing — `docker build --target test .` |
| `build-image` | `build` | `registry.granux.fr/decourt-demaria/tp/frontend` |
| `deploy` | `deploy` | triggers the CD pipeline of the `ci-cd` project |

Each build is pushed as `:<short SHA>`, plus `:latest` on the default branch and `:<git tag>` on a
tag, so a deployment can always name the exact commit it runs. `docker build --provenance=false` is
required: since Docker 29.7 the attestation can reach the registry before the image it refers to,
which GitLab rejects with `blob unknown to registry`.

The `deploy` job is a `trigger` job: it starts the pipeline of the
[`ci-cd`](../ci-cd) project with
`FRONTEND_TAG: $CI_COMMIT_SHORT_SHA`, so the CD deploys the image this pipeline just built, and
`strategy: depend` makes this job red if the deployment fails.
