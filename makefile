# Include environment variables from .env file
-include .env
export

BINARY_NAME=app

BUILD_DIR=build

build:
	@echo "==> Building binary..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags "-w -s" ./app
	@echo "==> Build complete! Binary is located at $(BUILD_DIR)/$(BINARY_NAME)"

run: build
	@echo "==> Running the application..."
	./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	@echo "==> Cleaning up..."
	rm -rf $(BUILD_DIR)
	@echo "==> Clean complete!"

test:
	@echo "==> Running tests..."
	go test ./...
	@echo "==> Tests complete!"

test-coverage:
	@echo "==> Running tests with coverage..."
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
	@echo "==> Coverage report generated!"

deps:
	@echo "==> Updating dependencies..."
	go mod tidy
	go mod download
	@echo "==> Dependencies updated!"

.PHONY: default
default: build

.PHONY: build run clean test test-coverage deps
