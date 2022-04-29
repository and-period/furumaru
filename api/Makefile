.PHONY: setup install mockgen protoc fmt vet lint test build

GOLANGCI_VERSION := 1.45.2

LINT_PACKAGES := $(shell go list $(CURDIR)/... | grep -v -e "mock" -v -e "proto" -v -e "third_party" -v -e "tmp")
TEST_PACKAGES := $(shell go list $(CURDIR)/internal/... $(CURDIR)/pkg/...)
EXTERNAL_PACKAGE=$(CURDIR)/third_party

setup: install
	wget -O - -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${GOLANGCI_VERSION}

install:
	go install github.com/golang/mock/mockgen@v1.6.0

mockgen:
	rm -rf ./mock
	go generate ./...

fmt:
	! gofmt -d -s ./cmd ./config ./hack ./internal ./pkg | grep '^'

vet:
	go vet $(LINT_PACKAGES)

lint:
	./bin/golangci-lint run -c .golangci.yaml ./...

test:
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic $(TEST_PACKAGES)

build:
	go build -o ./app ./cmd/${SERVICE}/main.go

dev:
	go build -o ./cmd/${SERVICE}/app ./cmd/${SERVICE}/main.go
	./cmd/${SERVICE}/app
