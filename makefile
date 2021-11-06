# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

SERVER_BINARY_NAME=fpsmonitor_server
SERVER_BINARY_UNIX=$(SERVER_BINARY_NAME)_unix
CLIENT_BINARY_NAME=fpsmonitor_client
CLIENT_BINARY_UNIX=$(CLIENT_BINARY_NAME)_unix

.PHONY: build-server build-client

all: build-server build-client

mkdir:
	mkdir -p ./build

build-server: mkdir
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/$(SERVER_BINARY_NAME) -v ./cmd/server

build-client: mkdir
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/$(CLIENT_BINARY_NAME) -v ./cmd/client

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f ./build/$(SERVER_BINARY_NAME)
	rm -f ./build/$(CLIENT_BINARY_NAME)

run-server: build-server
	cd ./build && ./$(SERVER_BINARY_NAME)

run-client: build-client
	cd ./build && ./$(CLIENT_BINARY_NAME)

tidy:
	$(GOMOD) tidy
