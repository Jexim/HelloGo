.PHONY: all swagger build run run-mock

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
