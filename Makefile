# Makefile for the Konsume project

BINARY_NAME=konsume
GOBUILD=go build
GOCLEAN=go clean
GOTEST=go test
GOGET=go get
GORUN=go run

.PHONY: all build test clean run deps plugin_postgres plugin_mongodb

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v .

test:
	$(GOTEST) -v ./... -race

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GORUN) .

deps:
	$(GOGET) ./...

plugin_postgres:
	$(GOBUILD) -buildmode=plugin -o ./plugins/postgres-darwin.so ./plugin/postgresql/postgresql.go

plugin_mongodb:
	$(GOBUILD) -buildmode=plugin -o ./plugins/mongodb-darwin.so ./plugin/mongodb/mongodb.go

start: plugin_postgres plugin_mongodb run