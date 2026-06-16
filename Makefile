# Common developer tasks. All commands assume you're in a `nix develop` shell
# or have the equivalent tools on $PATH.

.PHONY: help run test lint fmt build sqlc migrate-up migrate-down migrate-new openapi docker-up docker-down

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "Available targets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-18s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

run: ## Run the server locally
	go run ./cmd/server

test: ## Run all tests with race detector
	go test -race -cover ./...

lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format all Go files
	gofmt -w .
	go mod tidy

build: ## Build the server binary
	go build -o bin/server ./cmd/server

sqlc: ## Regenerate sqlc code from queries.sql files
	sqlc generate

migrate-up: ## Apply all pending migrations
	goose -dir migrations postgres "$$DATABASE_URL" up

migrate-down: ## Roll back one migration
	goose -dir migrations postgres "$$DATABASE_URL" down

migrate-new: ## Create a new migration. Usage: make migrate-new name=add_sessions
	goose -dir migrations create $(name) sql

openapi: ## Generate Go server interfaces from openapi.yaml
	oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml

docker-up: ## Start the local dev stack (app + postgres)
	docker compose -f deploy/docker/compose.yaml up --build

docker-down: ## Stop and remove the dev stack
	docker compose -f deploy/docker/compose.yaml down
