.PHONY: help build run migrate migrate-down migrate-status migrate-create clean install-deps dev

# Database URL for migrations (set from .env or override here)
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

help:
	@echo "Available commands:"
	@echo "  make install-deps      - Install Go dependencies"
	@echo "  make build             - Build the application"
	@echo "  make run               - Run the application"
	@echo "  make migrate           - Run all pending migrations"
	@echo "  make migrate-down      - Rollback last migration"
	@echo "  make migrate-status    - Show migration version"
	@echo "  make migrate-create    - Create a new migration file"
	@echo "  make dev               - Run with hot reload (requires air)"
	@echo "  make clean             - Clean build artifacts"

install-deps:
	go mod download
	go mod tidy

build:
	go build -o bin/api cmd/api/main.go

run: build
	./bin/api

migrate:
	migrate -path db/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DATABASE_URL)" down

migrate-status:
	migrate -path db/migrations -database "$(DATABASE_URL)" version

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide migration name"; \
		echo "Usage: make migrate-create name=<migration_name>"; \
		exit 1; \
	fi
	migrate create -ext sql -dir db/migrations -seq $(name)

dev:
	air

clean:
	rm -f bin/api
	rm -rf tmp/

.env:
	cp .env.example .env
	@echo "✓ .env file created. Please update with your configuration."
