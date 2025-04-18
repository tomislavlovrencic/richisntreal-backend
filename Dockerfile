# Stage 1: build the Go binary
FROM golang:1.24-alpine AS builder

# Install git+ca‑certs for dependencies, set working dir
RUN apk add --no-cache git ca-certificates
WORKDIR /app

# Copy go.mod / go.sum and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your code & build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o richisntreal cmd/main.go

# Stage 2: minimal runtime image
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary and migrations (so runMigrations can see them)
COPY --from=builder /app/richisntreal .
COPY --from=builder /app/internal/infrastructure/mysql/migrations ./migrations

# Expose your HTTP port
EXPOSE 8080

# Run the service
ENTRYPOINT ["./richisntreal"]
