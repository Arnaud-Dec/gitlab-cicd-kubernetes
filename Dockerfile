FROM golang:1.27-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Checks du README, lancés avec `docker build --target lint .` / `--target test .`
# Un build normal (image finale, docker compose) ne passe pas par ces étapes.
FROM base AS lint
# Cache de compilation conservé entre deux builds : sinon golangci-lint est recompilé à chaque fois
RUN --mount=type=cache,target=/root/.cache go vet ./...
RUN --mount=type=cache,target=/root/.cache go tool golangci-lint run ./...

FROM base AS test
RUN go test ./...
# Relancé contre un vrai Postgres par compose.test.yml (via TEST_DATABASE_DSN)
CMD ["go", "test", "./..."]

FROM base AS builder
RUN CGO_ENABLED=0 go build -ldflags "-X todoapp/internal/version.Version=$(date -u +%Y-%m-%dT%H:%M)" -o bin/api ./cmd/api

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/bin/api /app/api
EXPOSE 8080 9090
ENTRYPOINT ["/app/api"]
