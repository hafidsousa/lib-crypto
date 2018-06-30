# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOLINT=golangci-lint
BINARY_NAME=lib-crypto

all: lint test build

clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
lint:
		$(GOLINT) run --config .golangci.yml
test:
		$(GOTEST) ./... -short
build:
		$(GOBUILD) -o $(BINARY_NAME) -v
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)

# Cross compilation
build-linux: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v

.PHONY: all lint test build run