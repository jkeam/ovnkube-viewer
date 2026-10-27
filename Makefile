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

frontend-build: ## Build frontend for production (standalone)
	cd frontend && npm run build:standalone

frontend-build-plugin: ## Build frontend as console plugin
	cd frontend && npm run build:plugin

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

# OpenShift/Kubernetes deployment targets
REGISTRY ?= quay.io/yourorg
BACKEND_IMAGE ?= $(REGISTRY)/ovnkube-viewer-backend
FRONTEND_IMAGE ?= $(REGISTRY)/ovnkube-viewer-frontend
TAG ?= latest
OVERLAY ?= dev

.PHONY: openshift-build openshift-push openshift-deploy openshift-undeploy openshift-logs openshift-validate

openshift-build: ## Build container images for OpenShift
	podman build -t $(BACKEND_IMAGE):$(TAG) -f backend/Dockerfile backend
	podman build -t $(FRONTEND_IMAGE):$(TAG) -f frontend/Dockerfile frontend
	podman build -t $(REGISTRY)/ovnkube-viewer-plugin:$(TAG) -f frontend/Dockerfile.plugin frontend

openshift-push: ## Push container images to registry
	podman push $(BACKEND_IMAGE):$(TAG)
	podman push $(FRONTEND_IMAGE):$(TAG)
	podman push $(REGISTRY)/ovnkube-viewer-plugin:$(TAG)

openshift-validate: ## Validate Kustomize manifests
	@echo "Validating $(OVERLAY) overlay..."
	oc kustomize deploy/openshift/overlays/$(OVERLAY)

openshift-deploy: ## Deploy to OpenShift (use OVERLAY=dev or OVERLAY=prod)
	@echo "Deploying $(OVERLAY) overlay..."
	oc apply -k deploy/openshift/overlays/$(OVERLAY)
	@echo "Waiting for deployments to be ready..."
	@sleep 5
	@if [ "$(OVERLAY)" = "dev" ]; then \
		oc wait --for=condition=available --timeout=300s deployment/dev-ovnkube-viewer-backend -n ovnkube-viewer-dev || true; \
		oc wait --for=condition=available --timeout=300s deployment/dev-ovnkube-viewer-frontend -n ovnkube-viewer-dev || true; \
		echo "Getting route URL..."; \
		oc get route -n ovnkube-viewer-dev; \
	else \
		oc wait --for=condition=available --timeout=300s deployment/ovnkube-viewer-backend -n ovnkube-viewer || true; \
		oc wait --for=condition=available --timeout=300s deployment/ovnkube-viewer-frontend -n ovnkube-viewer || true; \
		echo "Getting route URL..."; \
		oc get route -n ovnkube-viewer; \
	fi

openshift-undeploy: ## Remove deployment from OpenShift
	@echo "Removing $(OVERLAY) overlay..."
	oc delete -k deploy/openshift/overlays/$(OVERLAY)

openshift-logs: ## Tail logs from OpenShift deployment
	@if [ "$(OVERLAY)" = "dev" ]; then \
		oc logs -f deployment/dev-ovnkube-viewer-backend -n ovnkube-viewer-dev; \
	else \
		oc logs -f deployment/ovnkube-viewer-backend -n ovnkube-viewer; \
	fi

openshift-status: ## Check OpenShift deployment status
	@if [ "$(OVERLAY)" = "dev" ]; then \
		echo "=== Pods ==="; \
		oc get pods -n ovnkube-viewer-dev -l app=ovnkube-viewer; \
		echo ""; \
		echo "=== Deployments ==="; \
		oc get deployments -n ovnkube-viewer-dev -l app=ovnkube-viewer; \
		echo ""; \
		echo "=== Services ==="; \
		oc get services -n ovnkube-viewer-dev -l app=ovnkube-viewer; \
		echo ""; \
		echo "=== Routes ==="; \
		oc get routes -n ovnkube-viewer-dev; \
	else \
		echo "=== Pods ==="; \
		oc get pods -n ovnkube-viewer -l app=ovnkube-viewer; \
		echo ""; \
		echo "=== Deployments ==="; \
		oc get deployments -n ovnkube-viewer -l app=ovnkube-viewer; \
		echo ""; \
		echo "=== Services ==="; \
		oc get services -n ovnkube-viewer -l app=ovnkube-viewer; \
		echo ""; \
		echo "=== Routes ==="; \
		oc get routes -n ovnkube-viewer; \
	fi

openshift-all: openshift-build openshift-push openshift-deploy ## Build, push, and deploy to OpenShift
