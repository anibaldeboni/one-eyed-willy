BIN="./bin"
SRC=$(shell find . -name "*.go")
export APP=one-eyed-willy
export LDFLAGS="-w -s"

ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH), run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

.PHONY: fmt lint test install_deps clean

default: all

all: fmt test

fmt:
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	golangci-lint run -v

test: install_deps
	go test -v ./...

setup:
	go install github.com/swaggo/swag/cmd/swag@latest

docs:
	swag i --parseInternal --dir ./cmd/oew/,./internal/handler/,./pkg/utils

install_deps:
	go get -v ./...

run:
	go run ./cmd/oew/

build:
	go build -v -o ./bin/$(APP) ./cmd/oew

build-static:
	CGO_ENABLED=0 go build -race -v -o $(APP) -a -installsuffix cgo -ldflags $(LDFLAGS) .

clean:
	rm -rf $(BIN)
