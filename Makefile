.PHONY: build run test clean

# Build the application
build:
	go build -o bin/richisntreal-backend

# Run the application
run:
	go run main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy

# Generate database migrations
migrate:
	go run cmd/migrate/main.go

# Run linter
lint:
	golangci-lint run 