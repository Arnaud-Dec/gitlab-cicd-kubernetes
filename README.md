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
