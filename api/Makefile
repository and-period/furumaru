.PHONY: help setup install mockgen protoc fmt vet lint test build

AIR_VERSION := 1.50.0
GOFUMPT_VERSION := 0.6.0
GOLANGCI_VERSION := 1.56.2
MOCKGEN_VERSION := 0.4.0

LINT_PACKAGES := $(shell go list $(CURDIR)/... | grep -v -e "mock" -v -e "tmp")
TEST_PACKAGES := $(shell go list $(CURDIR)/internal/... $(CURDIR)/pkg/...)

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: install ## 初回環境構築用
	wget -O - -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${GOLANGCI_VERSION}

install: ## 依存ライブラリのインストール
	go install mvdan.cc/gofumpt@v${GOFUMPT_VERSION}
	go install go.uber.org/mock/mockgen@v${MOCKGEN_VERSION}
	go install github.com/cosmtrek/air@v${AIR_VERSION}

mockgen: ## ユニットテストで使用するモックの生成
	rm -rf ./mock
	go generate ./...

fmt: ## フォーマットが正しくない箇所の出力
	! gofumpt -d ./cmd ./config ./hack ./internal ./pkg | grep '^'

vet: ## コードの静的解析
	go vet $(LINT_PACKAGES)

lint: ## Linterの実行
	./bin/golangci-lint run -c .golangci.yaml ./...

test: ## テストの実行
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic $(TEST_PACKAGES)

build: ## アプリケーションのコンパイル
	GOOS=linux GOARCH=amd64 go build -o ./app ./cmd/${SERVICE}/main.go

build-dev: ## アプリケーションのコンパイル(dev環境用)
	GOOS=linux GOARCH=amd64 go build -o ./cmd/${SERVICE}/app ./cmd/${SERVICE}/main.go

start-dev: ## アプリケーションの起動(dev環境用)
	air -c ./config/${SERVICE}/.air.toml
