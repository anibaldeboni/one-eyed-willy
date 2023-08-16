BIN="./bin"
SRC=$(shell find . -name "*.go")
export APP=one-eyed-willy
export LDFLAGS="-w -s"

.PHONY: lint test install_deps install_linter install_swag clean watch install_air

default: all

all: lint test

lint: install_linter
	@golangci-lint run -v

test: install_deps
	@go test -v -race ./...

docs: install_swag
	swag i --parseInternal --dir ./cmd/oew/,./internal/handler/,./pkg/utils

install_deps:
	@go get -v ./...

run:
	@go run ./cmd/oew/

build:
	@go build -v -o ./bin/$(APP) ./cmd/oew

build-static:
	CGO_ENABLED=0 go build -race -v -o $(APP) -a -installsuffix cgo -ldflags $(LDFLAGS) .

clean:
	@rm -rf $(BIN)

watch: install_air
	@air

install_linter:
	@if [ "$(shell which golangci-lint)" = "" ]; then\
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.3;\
	fi
	
install_swag:
	@if [ "$(shell which swag)" = "" ]; then\
		go install github.com/swaggo/swag/cmd/swag@latest;\
	fi

install_air:
	@if [ "$(shell which air)" = "" ]; then\
		go install github.com/cosmtrek/air@latest;\
	fi
