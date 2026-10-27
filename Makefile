.PHONY: help backend-build frontend-build backend-run frontend-dev docker-build docker-up docker-down clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Backend targets
backend-build: ## Build backend server
	cd backend && go build -o bin/server ./cmd/server

backend-run: ## Run backend server
	cd backend && go run ./cmd/server

backend-test: ## Run backend tests
	cd backend && go test ./...

# Frontend targets
frontend-install: ## Install frontend dependencies
	cd frontend && npm install

frontend-dev: ## Run frontend in development mode
	cd frontend && npm run dev

frontend-build: ## Build frontend for production
	cd frontend && npm run build

# Docker targets
docker-build: ## Build Docker images
	cd deploy/standalone && docker-compose build

docker-up: ## Start services with docker-compose
	cd deploy/standalone && docker-compose up -d

docker-down: ## Stop services
	cd deploy/standalone && docker-compose down

docker-logs: ## View docker logs
	cd deploy/standalone && docker-compose logs -f

# Development
dev: ## Run both backend and frontend in dev mode (requires tmux or run in separate terminals)
	@echo "Run 'make backend-run' in one terminal and 'make frontend-dev' in another"

# Clean
clean: ## Clean build artifacts
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	cd backend && go clean

# Initialize
init: ## Initialize project (install deps, etc)
	cd backend && go mod tidy
	cd frontend && npm install

# Testing
test: backend-test ## Run all tests

# All
all: backend-build frontend-build ## Build everything
