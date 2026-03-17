BINARY        = server
BUILD_DIR     = ./cmd/api
MAIN          = $(BUILD_DIR)/main.go
MIGRATE       = migrate
MIGRATE_PATH  = migrations
DB_URL       ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

-include .env
export

.PHONY: build run dev test tidy docker-up docker-down clean swagger migrate-up migrate-down migrate-force

## build: compile the binary
build:
	go build -ldflags="-s -w" -o $(BINARY) $(BUILD_DIR)

## run: build & run the server
run: build
	./$(BINARY)

## dev: run with hot-reload using air (go install github.com/air-verse/air@latest)
dev:
	air

## test: run all tests
test:
	go test -v -race ./...

## tidy: tidy & vendor modules
tidy:
	go mod tidy

## docker-up: start all containers (postgres + api)
docker-up:
	docker compose up -d --build

## docker-down: stop all containers
docker-down:
	docker compose down

## swagger: regenerate Swagger docs (requires: go install github.com/swaggo/swag/cmd/swag@latest)
swagger:
	~/go/bin/swag init --generalInfo main.go --output docs

## migrate-up: apply all pending migrations
migrate-up:
	$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DB_URL)" up

## migrate-down: roll back the last migration
migrate-down:
	$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

## migrate-force: force set migration version (usage: make migrate-force V=3)
migrate-force:
	$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DB_URL)" force $(V)

## clean: remove build artifacts
clean:
	rm -f $(BINARY)

## help: display available targets
help:
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/## //'
