.PHONY: help install mockgen fmt vet lint lint-fix test build build-hack build-dev start-dev

BUILD_DIR := $(CURDIR)

LINT_PACKAGES := $(shell go list $(CURDIR)/... | grep -v -e "mock" -v -e "tmp")
TEST_PACKAGES := $(shell go list $(CURDIR)/internal/... $(CURDIR)/pkg/...)

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## 依存ライブラリのインストール
	go mod tidy
	go get tool

mockgen: ## ユニットテストで使用するモックの生成
	rm -rf ./mock
	go generate ./...

fmt: ## フォーマットが正しくない箇所の出力
	! go tool gofumpt -d ./cmd ./config ./hack ./internal ./pkg | grep '^'

vet: ## コードの静的解析
	go vet $(LINT_PACKAGES)

lint: ## Linterの実行
	go tool golangci-lint run -c .golangci.yaml ./...

lint-fix: ## Linterでの修正
	go tool golangci-lint run -c .golangci.yaml ./... --fix

test: ## テストの実行
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic $(TEST_PACKAGES)

build: ## アプリケーションのコンパイル
	GOOS=linux GOARCH=arm64 go build -o ${BUILD_DIR}/app ./cmd/${SERVICE}/main.go

build-hack: ## コマンドのコンパイル
	GOOS=linux GOARCH=arm64 go build -o ${BUILD_DIR}/app ./hack/${SERVICE}/main.go

build-dev: ## アプリケーションのコンパイル(dev環境用)
	GOOS=linux GOARCH=arm64 go build -o ./cmd/${SERVICE}/app ./cmd/${SERVICE}/main.go

start-dev: ## アプリケーションの起動(dev環境用)
	go tool air -c ./config/${SERVICE}/.air.toml
