# Makefile for the Konsume project

BINARY_NAME=konsume
GOBUILD=go build
GOCLEAN=go clean
GOTEST=go test
GOGET=go get
GORUN=go run

.PHONY: all build test clean run deps

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GORUN) .

deps:
	$(GOGET) ./...
