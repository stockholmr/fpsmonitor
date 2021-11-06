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

all: build

build: 
		mkdir ./build
		CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/$(SERVER_BINARY_NAME) -v ./cmd/server
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -R -f ./build

run:
		mkdir ./build
		$(GOBUILD) -o ./build/$(SERVER_BINARY_NAME) -v ./cmd/server
		./build/$(SERVER_BINARY_NAME)
tidy:
		$(GOMOD) tidy

# Cross compilation
build-windows:
		 $(GOBUILD) -o $(BINARY_UNIX) -v