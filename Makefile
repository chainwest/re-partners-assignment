.PHONY: help run build test clean lint fmt vet tidy docker-build docker-run compose-up compose-down compose-logs migrate smoke

APP_NAME=api
VERSION?=dev
PORT?=8080
BUILD_DIR=./build
CMD_DIR=./cmd/api
COMPOSE_FILE=./deployments/docker-compose.yaml

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application locally
	@PORT=$(PORT) VERSION=$(VERSION) go run $(CMD_DIR)/main.go

build: ## Build the application binary
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

test: ## Run all tests with coverage
	@go test -v -race -coverprofile=coverage.out ./...

coverage: test ## Show test coverage in browser
	@go tool cover -html=coverage.out

clean: ## Clean build artifacts
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out usecase_coverage.out

lint: ## Run linter (golangci-lint)
	@golangci-lint run ./... || echo "golangci-lint not installed. Run: make install-tools"

fmt: ## Format code with gofmt
	@go fmt ./...

vet: ## Run go vet
	@go vet ./...

tidy: ## Tidy and verify go modules
	@go mod tidy
	@go mod verify

docker-build: ## Build docker image
	@docker build -t $(APP_NAME):$(VERSION) -f deployments/Dockerfile .

docker-run: ## Run docker container standalone
	@docker run -p $(PORT):8080 --name $(APP_NAME) --rm $(APP_NAME):$(VERSION)

compose-up: ## Start all services (api, postgres, redis)
	@docker-compose -f $(COMPOSE_FILE) up -d

compose-down: ## Stop and remove all services
	@docker-compose -f $(COMPOSE_FILE) down

compose-logs: ## Show logs from all services
	@docker-compose -f $(COMPOSE_FILE) logs -f

compose-restart: ## Restart all services
	@docker-compose -f $(COMPOSE_FILE) restart

compose-build: ## Rebuild docker images
	@docker-compose -f $(COMPOSE_FILE) build --no-cache

compose-ps: ## Show status of services
	@docker-compose -f $(COMPOSE_FILE) ps

migrate: ## Run database migrations (requires postgres running)
	@docker-compose -f $(COMPOSE_FILE) up -d postgres
	@sleep 5
	@docker exec re-partners-postgres psql -U postgres -d re_partners -f /docker-entrypoint-initdb.d/001_create_pack_sets.up.sql || true
	@docker exec re-partners-postgres psql -U postgres -d re_partners -f /docker-entrypoint-initdb.d/002_create_calculations.up.sql || true

migrate-down: ## Rollback database migrations
	@docker exec re-partners-postgres psql -U postgres -d re_partners -f /docker-entrypoint-initdb.d/002_create_calculations.down.sql || true
	@docker exec re-partners-postgres psql -U postgres -d re_partners -f /docker-entrypoint-initdb.d/001_create_pack_sets.down.sql || true

smoke: ## Run smoke tests against running service
	@./scripts/smoke_test.sh

test-api: ## Test API endpoints (requires running server)
	@./scripts/test_api.sh

install-tools: ## Install development tools
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

dev: fmt vet run ## Format, vet and run the application

all: clean lint test build ## Run all checks and build

