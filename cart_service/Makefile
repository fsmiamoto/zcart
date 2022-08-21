GO := $(shell which go)
COVEROUT := coverage.out

all: build

build:
	$(GO) build -o ./bin/cart_service ./cmd/cart_service/

test:
	$(GO) test -coverprofile=$(COVEROUT) ./... -v

coverage: test
	$(GO) tool cover -html=$(COVEROUT)

dev:
	DEV_MODE=true CompileDaemon -build "make build" -command "./bin/cart_service"

.PHONY: all build test coverage
