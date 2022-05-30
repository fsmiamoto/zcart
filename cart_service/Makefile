GO := $(shell which go)

all: build

build:
	$(GO) build -o ./bin/cart_service ./cmd/cart_service/

dev:
	CompileDaemon -build "make build" -command "./bin/cart_service"

.PHONY: all build
