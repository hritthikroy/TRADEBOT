.PHONY: help build run test clean docker-build docker-run lint fmt

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building..."
	cd backend && go build -o trading-bot .
	@echo "✅ Build complete: backend/trading-bot"

run: ## Run the application
	@echo "Starting server..."
	cd backend && ./trading-bot

dev: ## Run in development mode with auto-reload
	@echo "Starting development server..."
	cd backend && go run .

test: ## Run tests
	@echo "Running tests..."
	cd backend && go test -v -race -coverprofile=coverage.out ./...
	@echo "✅ Tests complete"

test-coverage: test ## Run tests with coverage report
	cd backend && go tool cover -html=coverage.out

lint: ## Run linter
	@echo "Running linter..."
	cd backend && golangci-lint run
	@echo "✅ Lint complete"

fmt: ## Format code
	@echo "Formatting code..."
	cd backend && go fmt ./...
	@echo "✅ Format complete"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f backend/trading-bot
	rm -f backend/coverage.out
	rm -f backend/*.log
	@echo "✅ Clean complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	cd backend && docker build -t trading-bot:latest .
	@echo "✅ Docker build complete"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file backend/.env trading-bot:latest

docker-stop: ## Stop Docker container
	docker stop $$(docker ps -q --filter ancestor=trading-bot:latest)

install-deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	cd backend && go mod download
	@echo "✅ Dependencies installed"

install-tools: ## Install development tools
	@echo "Installing tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Tools installed"

backtest: ## Run a quick backtest
	@echo "Running backtest..."
	curl -X POST http://localhost:8080/api/v1/backtest/run \
		-H "Content-Type: application/json" \
		-d '{"symbol":"BTCUSDT","interval":"15m","days":30,"startBalance":500}'

health: ## Check server health
	@curl -s http://localhost:8080/api/v1/health | python3 -m json.tool

all: clean fmt lint test build ## Run all checks and build

.DEFAULT_GOAL := help
