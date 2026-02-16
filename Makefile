.PHONY: build run test test-short test-integration lint generate migrate-up migrate-down docker-up docker-down clean

# Build
build:
	go build -o bin/api ./cmd/api
	go build -o bin/batch ./cmd/batch

run:
	go run ./cmd/api

# Test
test:
	go test -race -count=1 ./...

test-short:
	go test -short -race ./...

test-integration:
	go test -race -count=1 ./internal/repository/...

# Lint
lint:
	go tool golangci-lint run

# Generate (mocks)
generate:
	go generate ./...

# Migration
migrate-up:
	go run ./cmd/batch migrate up

migrate-down:
	go run ./cmd/batch migrate down

# Docker
docker-up:
	docker compose -f devenv/compose.yml up -d

docker-down:
	docker compose -f devenv/compose.yml down

# Clean
clean:
	rm -rf bin/ tmp/
