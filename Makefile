.PHONY: build run test lint clean docker-up docker-down explore

# Build the server binary
build:
	go build -o bin/server ./cmd/server

# Run the server locally
run:
	go run ./cmd/server

# Run all tests
test:
	go test -v -race -count=1 ./...

# Run tests with coverage
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Lint (requires golangci-lint)
lint:
	golangci-lint run ./...

# Remove build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Start Redis via Docker Compose
docker-up:
	docker compose up -d

# Stop Docker Compose services
docker-down:
	docker compose down

# Run the explore CLI tool
explore:
	go run ./cmd/explore
