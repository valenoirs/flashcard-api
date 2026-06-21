# Loading the environment variable
include .env
export

MIGRATE := $(shell go env GOPATH)/bin/migrate
TEST_LAYER ?= service|handler
SCHEMA ?= 1
MAIN_PATH := ./cmd/
BUILD_PATH := ./bin/app

.PHONY: run build test lint format clean help
.DEFAULT_GOAL := help

run:
	go run $(MAIN_PATH)

build:
	go build -ldflags="-s -w" -o $(BUILD_PATH) $(MAIN_PATH)

test:
	go test ./...

lint:
	go vet ./...

format:
	go fmt ./...

clean:
	rm -rf bin/ coverage.out

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

coverage:
	@echo "Gathering packages matching: $(LAYER)..."
	@PACKAGES=$$(go list ./internal/... | grep -E '/($(TEST_LAYER))' || true); \
	if [ -z "$$PACKAGES" ]; then \
		echo "No matching packages found for pattern: $(TEST_LAYER)"; \
		exit 1; \
	fi; \
	echo "Running tests on targeted packages..."; \
	go test -coverprofile=coverage.out $$PACKAGES; \

coverage-show: coverage
	@echo "Coverage Summary:"
	@go tool cover -func=coverage.out | grep total:
	@go tool cover -html=coverage.out
	@echo "HTML report written to coverage.html"

create-migration:
	@if [ -z "$(NAME)" ]; then echo "Error: NAME is required. Usage: make create-migration SCHEMA=1 NAME=create_user"; exit 1; fi
	@mkdir -p migration
	@DATE=$$(date +%Y%m%d); \
	LATEST=$$(ls migration/$(SCHEMA)[0-9][0-9][0-9]_*.sql 2>/dev/null | awk -F'/' '{print $$NF}' | cut -c 2-4 | sort -n | tail -1); \
	if [ -z "$$LATEST" ]; then \
		NEXT="001"; \
	else \
		NEXT=$$(echo $$LATEST | awk '{printf "%03d", $$1 + 1}'); \
	fi; \
	FILENAME="$(SCHEMA)$${NEXT}_$${DATE}_$(NAME)"; \
	echo "Creating migration:"; \
	echo "  -> migration/$$FILENAME.up.sql"; \
	echo "  -> migration/$$FILENAME.down.sql"; \
	touch "migration/$$FILENAME.up.sql"; \
	touch "migration/$$FILENAME.down.sql"

migration-up:
	@echo "Running migration up..."
	@echo "Initializing schema..."
	@$(MIGRATE) -path migration/init -database "$(DB_CONNECTION_STRING)" up
	@echo "Executing migration..."
	@$(MIGRATE) -path migration -database "$(DB_CONNECTION_STRING)&search_path=$(DB_SCHEMA_NAME)" up

migration-down:
	@echo "Running migration down..."
	@$(MIGRATE) -path migration -database "$(DB_CONNECTION_STRING)&search_path=$(DB_SCHEMA_NAME)" down 1

reset-migration:
	@echo "Force reset migration to $(VERSION)"
	@$(MIGRATE) -path migration -database "$(DB_CONNECTION_STRING)&search_path=$(DB_SCHEMA_NAME)" force $(VERSION)

migration-version:
	@echo "Current migration version"
	@$(MIGRATE) -path migration -database "$(DB_CONNECTION_STRING)&search_path=$(DB_SCHEMA_NAME)" version

create-seeder:
	@if [ -z "$(NAME)" ]; then echo "Error: NAME is required. Usage: make create-seeder SCHEMA=1 NAME=create_user"; exit 1; fi
	@mkdir -p seeder
	@DATE=$$(date +%Y%m%d); \
	LATEST=$$(ls seeder/$(SCHEMA)[0-9][0-9][0-9]_*.sql 2>/dev/null | awk -F'/' '{print $$NF}' | cut -c 2-4 | sort -n | tail -1); \
	if [ -z "$$LATEST" ]; then \
		NEXT="001"; \
	else \
		NEXT=$$(echo $$LATEST | awk '{printf "%03d", $$1 + 1}'); \
	fi; \
	FILENAME="$(SCHEMA)$${NEXT}_$${DATE}_$(NAME)"; \
	echo "Creating seeder:"; \
	echo "  -> seeder/$$FILENAME.up.sql"; \
	touch "seeder/$$FILENAME.up.sql"; \

seeder-up:
	@echo "Running seeder up..."
	@$(MIGRATE) -path seeder -database "$(DB_CONNECTION_STRING)&search_path=$(DB_SCHEMA_NAME)&x-migrations-table=schema_seeders" up
