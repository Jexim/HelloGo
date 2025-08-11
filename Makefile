.PHONY: all swagger build run run-mock goose-install migrate-up migrate-down migrate-status migrate-create sqlc

all: build

swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --generalInfo cmd/server/main.go --output ./docs

build:
	go build -o bin/server ./cmd/server

# Default port if not specified
PORT ?= 8080

run:
	SERVER_ADDRESS=:$(PORT) go run cmd/server/main.go

run-mock:
	cd mocks/camera && go run main.go -port $(PORT)

# sqlc codegen
sqlc:
	sqlc generate -f ./db/sqlc.yaml

# goose migrations
goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

# Set DB_DSN env var before running (e.g., export DB_DSN="host=localhost user=postgres password=postgres dbname=hello port=5432 sslmode=disable")
migrate-up:
	goose -dir db/migrations postgres "$$DB_DSN" up

migrate-down:
	goose -dir db/migrations postgres "$$DB_DSN" down

migrate-status:
	goose -dir db/migrations postgres "$$DB_DSN" status

migrate-create:
	goose -dir db/migrations create $(name) sql
