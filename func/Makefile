.PHONY: help setup install create update fmt vet lint test build

GOFUMPT_VERSION := 0.4.0
GOLANGCI_VERSION := 1.52.2

LAMBDA_ROLE :=  arn:aws:iam::386661535629:role/service-role/furumaru-stg-lambda

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: install ## 初回環境構築用
	wget -O - -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${GOLANGCI_VERSION}

install: ## 依存ライブラリのインストール
	go install mvdan.cc/gofumpt@v${GOFUMPT_VERSION}

create: build ## Lambda関数の作成
	aws lambda create-function \
		--function-name ${FUNC_NAME} \
		--zip-file fileb://function.zip \
		--handler main \
		--runtime go1.x \
		--role ${LAMBDA_ROLE}

update: build ## Lambda関数のコードを更新
	aws lambda update-function-code \
		--function-name ${FUNC_NAME} \
		--zip-file fileb://function.zip

fmt: ## フォーマットが正しくない箇所の出力
	! gofumpt -d . | grep '^'

vet: ## コードの静的解析
	go vet ./...

lint: ## Linterの実行
	./bin/golangci-lint run -c .golangci.yaml ./...

test: ## テストの実行
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic ./...

build: ## アプリケーションのコンパイルと生成ファイルの圧縮
	GOOS=linux GOARCH=amd64 go build -o main ${FUNC_DIR}/main.go
	zip function.zip main
