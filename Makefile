.PHONY: build run test clean install-deps fmt lint docker-build docker-run

# Build variables
BINARY_NAME=bu-agent
BUILD_DIR=./bin
MAIN_PATH=./cmd/agent

# Build the application
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run: build
	@echo "Running..."
	@$(BUILD_DIR)/$(BINARY_NAME) -config config.yaml

# Run with debug logging
run-debug: build
	@echo "Running with debug logging..."
	@$(BUILD_DIR)/$(BINARY_NAME) -config config.yaml -log-level debug

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	@go tool cover -html=coverage.out

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@go vet ./...
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t bu-agent:latest .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run --rm -v $(PWD)/config.yaml:/app/config.yaml bu-agent:latest

#### Generate example config
###example-config:
###	@echo "Generating example config..."
###	@cat > config.example.yaml << EOF
#### Blockchain Universe Agent Configuration
###
###agent:
###  decision_interval: 30s
###  max_event_chain: 100
###
###llm:
###  api_endpoint: "http://localhost:11434/v1/completions"
###  api_key: ""
###  model: "llama3.2"
###  max_tokens: 150
###  temperature: 0.7
###  timeout_seconds: 30
###EOF
###	@echo "Example config generated: config.example.yaml"
###
# Help command
help:
	@echo "Available commands:"
	@echo "  make build           - Build the application"
	@echo "  make run             - Build and run the application"
	@echo "  make run-debug       - Run with debug logging"
	@echo "  make test            - Run tests"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make install-deps    - Install dependencies"
	@echo "  make fmt             - Format code"
	@echo "  make lint            - Lint code"
	@echo "  make docker-build    - Build Docker image"
	@echo "  make docker-run      - Run Docker container"
	@echo "  make example-config  - Generate example config"
