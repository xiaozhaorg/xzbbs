.PHONY: dev build clean run

# Development with hot reload (requires: go install github.com/air-verse/air@latest)
dev:
	air -c .air.toml

# Build the binary
build:
	cd web && npm run build
	go build -o bin/XzBBS ./cmd/server

# Build without frontend
build-api:
	go build -o bin/XzBBS ./cmd/server

# Run
run:
	go run ./cmd/server

# Clean
clean:
	rm -rf bin/ tmp/

# Docker
docker:
	docker build -t XzBBS .

# Generate go.sum
tidy:
	go mod tidy
