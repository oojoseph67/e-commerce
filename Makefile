.PHONY: help build run dev lint format migrate-up migrate-down migrate-force migrate-alter docker-up docker-down docs-generate

GOROOT=/opt/homebrew/opt/go/libexec
export GOROOT

.DEFAULT_GOAL := help

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -ldflags="-s -w" -o bin/app ./cmd/api

run: ## Run the application
	go run ./cmd/api

dev: ## Run the application with live reload
	go run ./cmd/api

lint: ## Run linter
	golangci-lint run ./...

format: ## Run the go formatter to format code and rearrange imports
	@gofmt -s -w .
	@$(shell go env GOPATH)/bin/goimports -w .

migrate-up: ## Run database migrations up
	migrate -path ./internal/database/migrations -database "postgresql://postgres:admin@localhost:5433/go-ecommerce?sslmode=disable" up

migrate-down: ## Run database migrations down
	migrate -path ./internal/database/migrations -database "postgresql://postgres:admin@localhost:5433/go-ecommerce?sslmode=disable" down

migrate-force: ## Run database migration force (usage: make migrate-alter version=1)
	migrate -path ./internal/database/migrations -database "postgresql://postgres:admin@localhost:5433/go-ecommerce?sslmode=disable" force $(VERSION)

migrate-alter: ## Create a new migration file (usage: make migrate-alter NAME=add_slug_to_products)
	migrate create -ext sql -dir ./internal/database/migrations $(NAME)

docker-up: ## Run docker up
	docker compose -f docker/docker-compose.yml up -d

docker-down: ## Run docker down
	docker compose -f docker/docker-compose.yml down

docs-generate: ## Generate swagger generate
	mkdir -p docs
	$(shell go env GOPATH)/bin/swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --exclude .git,docs,docker,database
