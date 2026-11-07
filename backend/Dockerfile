FROM golang:1.27-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Checks du README, lancés avec `docker build --target lint .` / `--target test .`
# Un build normal (image finale, docker compose) ne passe pas par ces étapes.
FROM base AS lint
# Binaire précompilé : compiler golangci-lint via `go tool` prend plusieurs minutes sur le runner.
# Version à garder identique à la directive `tool` de go.mod.
COPY --from=golangci/golangci-lint:v2.13.2-alpine /usr/bin/golangci-lint /usr/bin/golangci-lint
RUN --mount=type=cache,target=/root/.cache go vet ./...
RUN --mount=type=cache,target=/root/.cache golangci-lint run ./...

FROM base AS test
RUN go test ./...
# Relancé contre un vrai Postgres par compose.test.yml (via TEST_DATABASE_DSN)
CMD ["go", "test", "./..."]

# Spec OpenAPI 3 générée depuis les annotations du code (README, § Generate the OpenAPI spec)
FROM base AS openapi
RUN --mount=type=cache,target=/root/.cache go tool swag init -g cmd/api/main.go -o docs --outputTypes json
# chmod : cmd/openapi écrit le fichier en 0600, illisible pour nginx dans l'image Swagger UI
RUN --mount=type=cache,target=/root/.cache go run ./cmd/openapi && chmod 644 docs/openapi.json

# Image Swagger UI qui sert la spec : docker build --target swagger .
FROM swaggerapi/swagger-ui:v5.32.15 AS swagger
COPY --from=openapi /app/docs/openapi.json /spec/openapi.json
ENV SWAGGER_JSON=/spec/openapi.json

FROM base AS builder
RUN CGO_ENABLED=0 go build -ldflags "-X todoapp/internal/version.Version=$(date -u +%Y-%m-%dT%H:%M)" -o bin/api ./cmd/api

# Image de l'API : dernière étape, donc construite par défaut (docker build ., docker compose)
FROM alpine:3.21 AS api
WORKDIR /app
COPY --from=builder /app/bin/api /app/api
EXPOSE 8080 9090
ENTRYPOINT ["/app/api"]
