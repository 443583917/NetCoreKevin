.PHONY: build run test test-unit test-integration test-e2e test-coverage lint clean docker-build docker-run docker-stop

# Variables
APP_NAME = go-kevin
BUILD_DIR = bin
MAIN_PATH = cmd/server/main.go

# Build
build:
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Run
run:
	go run $(MAIN_PATH)

# Test all
test:
	go test ./...

# Unit tests
test-unit:
	go test ./internal/... ./pkg/...

# Integration tests
test-integration:
	go test ./test/integration/... -v

# API tests
test-e2e:
	go test ./test/e2e/... -v

# Test coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

# Lint
lint:
	golangci-lint run ./...

# Clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Docker build
docker-build:
	docker build -t $(APP_NAME) .

# Docker run
docker-run:
	docker-compose up -d

# Docker stop
docker-stop:
	docker-compose down

# Generate docs
docs:
	swag init -g cmd/server/main.go -o docs/api
