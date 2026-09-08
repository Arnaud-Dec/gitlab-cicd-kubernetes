FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-X todoapp/internal/version.Version=$(date -u +%Y-%m-%dT%H:%M)" -o bin/api ./cmd/api

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/bin/api /app/api
EXPOSE 8080 9090
ENTRYPOINT ["/app/api"]
