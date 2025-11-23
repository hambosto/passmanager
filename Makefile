.PHONY: build test clean install run lint fmt deps help

APP_NAME=passmanager
VERSION=1.0.0
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin

build:
	@echo "Building $(APP_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w -X main.Version=$(VERSION)" \
		-o $(BUILD_DIR)/$(APP_NAME) cmd/passmanager/main.go
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 cmd/passmanager/main.go
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 cmd/passmanager/main.go
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 cmd/passmanager/main.go
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe cmd/passmanager/main.go
	@echo "Multi-platform build complete"

test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-short:
	@echo "Running tests (short)..."
	@go test -short ./...

bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

install: build
	@echo "Installing $(APP_NAME)..."
	@cp $(BUILD_DIR)/$(APP_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(APP_NAME)
	@echo "Installed to $(INSTALL_DIR)/$(APP_NAME)"

run: build
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

lint:
	@echo "Running linters..."
	@go vet ./...
	@echo "Linting complete"

fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .
	@echo "Formatting complete"

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  build-all    - Build for all platforms"
	@echo "  test         - Run tests with coverage"
	@echo "  test-short   - Run tests without coverage"
	@echo "  bench        - Run benchmarks"
	@echo "  clean        - Remove build artifacts"
	@echo "  install      - Install to system"
	@echo "  run          - Build and run the application"
	@echo "  lint         - Run linters"
	@echo "  fmt          - Format code"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  help         - Show this help message"
