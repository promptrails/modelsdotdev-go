.PHONY: all build test lint fmt vet coverage clean update-bundle

all: fmt vet lint test build

## Build all packages
build:
	go build ./...

## Run all tests
test:
	go test -race -count=1 ./...

## Run tests with coverage
coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

## Run linters
lint:
	golangci-lint run ./...

## Format code
fmt:
	gofmt -w .
	goimports -w .

## Run go vet
vet:
	go vet ./...

## Refresh the embedded offline snapshot from models.dev
update-bundle:
	curl -fsSL https://models.dev/api.json -o data/models.json
	@echo "bundle updated ($$(wc -c < data/models.json) bytes)"

## Clean build artifacts
clean:
	rm -f coverage.out coverage.html
