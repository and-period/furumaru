.PHONY: help install mockgen fmt fmt-fix lint lint-fix test build build-hack build-dev start-dev

BUILD_DIR := $(CURDIR)
TEST_PACKAGES := $(shell go list $(CURDIR)/internal/... $(CURDIR)/pkg/...)

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## 依存ライブラリのインストール
	go mod download
	go get tool

mockgen: ## ユニットテストで使用するモックの生成
	rm -rf ./mock
	go generate ./...

fmt:
	go tool golangci-lint fmt -c .golangci.yaml --diff ./pkg/... ./internal/...

fmt-fix:
	go tool golangci-lint fmt -c .golangci.yaml ./pkg/... ./internal/...

lint:
	go tool golangci-lint run -c .golangci.yaml ./...

lint-fix:
	go tool golangci-lint run -c .golangci.yaml ./... --fix

test: ## テストの実行
	go tool gotestsum -- -v -cover -coverprofile=coverage.txt -covermode=atomic $(TEST_PACKAGES)

build: ## アプリケーションのコンパイル
	GOOS=linux GOARCH=arm64 go build -o ${BUILD_DIR}/app ./cmd/${SERVICE}/main.go

build-hack: ## コマンドのコンパイル
	GOOS=linux GOARCH=arm64 go build -o ${BUILD_DIR}/app ./hack/${SERVICE}/main.go

build-dev: ## アプリケーションのコンパイル(dev環境用)
	GOOS=linux GOARCH=arm64 go build -o ./cmd/${SERVICE}/app ./cmd/${SERVICE}/main.go

start-dev: ## アプリケーションの起動(dev環境用)
	go tool air -c ./config/${SERVICE}/.air.toml
