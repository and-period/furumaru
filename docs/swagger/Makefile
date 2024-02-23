.PHONY: help setup build generate lint validate

OPENAPITOOLS_VERSION := 2.9.0
SWAGGER_CLI_VERSION  := 4.0.4
SWAGGER_UI_VERSION   := 5.11.8

ADMIN_INPUT_DIR       := admin
ADMIN_OUTPUT_DOCS_DIR := ./../../tmp/data/swagger/admin/docs
ADMIN_OUTPUT_TS_DIR   := ./../../tmp/data/swagger/admin/types
ADMIN_WEB_TYPES_DIR   := ./../../web/admin/src/types/api

USER_INPUT_DIR       := user
USER_OUTPUT_DOCS_DIR := ./../../tmp/data/swagger/user/docs
USER_OUTPUT_TS_DIR   := ./../../tmp/data/swagger/user/types
USER_WEB_TYPES_DIR   := ./../../web/user/src/types/api

ARTIFACTS_DIR := artifacts

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## 環境構築
	yarn global add \
		@openapitools/openapi-generator-cli@v${OPENAPITOOLS_VERSION} \
		swagger-cli@v${SWAGGER_CLI_VERSION}

build: build-admin build-user ## API仕様書の生成

build-admin: ## API仕様書の生成(管理者関連)
#	$(MAKE) generate TYPE=html2 INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=typescript-axios INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_TS_DIR}
	mkdir -p ${ADMIN_WEB_TYPES_DIR}; cp -r ${ADMIN_OUTPUT_TS_DIR}/*.ts ${ADMIN_WEB_TYPES_DIR}

build-user: ## API仕様書の生成(購入者関連)
#	$(MAKE) generate TYPE=html2 INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=typescript-fetch INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_TS_DIR}
	mkdir -p ${USER_WEB_TYPES_DIR}; cp -r ${USER_OUTPUT_TS_DIR}/*.ts ${USER_OUTPUT_TS_DIR}/apis ${USER_OUTPUT_TS_DIR}/models ${USER_WEB_TYPES_DIR}

build-docs: ## ホスティング環境へデプロイ用の成果物を生成
	$(MAKE) download-swagger-ui OUTPUT_DIR=${ARTIFACTS_DIR}
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${INPUT} OUTPUT_DIR=${ARTIFACTS_DIR}
	swagger-cli bundle ${ARTIFACTS_DIR}/openapi/openapi.yaml  --outfile ${ARTIFACTS_DIR}/swagger/openapi.yaml --type yaml --format 3
	sed -i 's/https:\/\/petstore.swagger.io\/v2\/swagger.json/openapi.yaml/' ${ARTIFACTS_DIR}/swagger/swagger-initializer.js

build-html: ## API仕様書の生成(HTML形式)
	$(MAKE) generate TYPE=html2 INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=html2 INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}

build-yaml: ## API仕様書の生成(OpenAPI形式)
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}

build-typescript: ## API仕様書の生成(TypeScript型定義とクライアント生成)
	$(MAKE) generate TYPE=typescript-axios INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_TS_DIR}
	$(MAKE) generate TYPE=typescript-fetch INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_TS_DIR}
	mkdir -p ${ADMIN_WEB_TYPES_DIR}; cp -r ${ADMIN_OUTPUT_TS_DIR}/*.ts ${ADMIN_WEB_TYPES_DIR}
	mkdir -p ${USER_WEB_TYPES_DIR}; cp -r ${USER_OUTPUT_TS_DIR}/*.ts ${USER_OUTPUT_TS_DIR}/apis ${USER_OUTPUT_TS_DIR}/models ${USER_WEB_TYPES_DIR}

generate: ## API仕様書の生成
	openapi-generator-cli generate \
		-g ${TYPE} \
		-i ${INPUT_DIR}/openapi.yaml \
		-o ${OUTPUT_DIR}

lint: ## OpenAPI形式に則っているかの構文チェック
	$(MAKE) validate INPUT_DIR=${ADMIN_INPUT_DIR}
	$(MAKE) validate INPUT_DIR=${USER_INPUT_DIR}

validate: ## OpenAPI形式に則っているかの構文チェック
	openapi-generator-cli validate \
		-i ${INPUT_DIR}/openapi.yaml \
		--recommend

download-swagger-ui: ## ホスティング時に使用するSwaggerUIのダウンロード
	mkdir -p artifacts/swagger
	wget https://github.com/swagger-api/swagger-ui/archive/refs/tags/v${SWAGGER_UI_VERSION}.tar.gz
	tar -zxvf v${SWAGGER_UI_VERSION}.tar.gz swagger-ui-${SWAGGER_UI_VERSION}/dist/
	cp swagger-ui-${SWAGGER_UI_VERSION}/dist/* ${OUTPUT_DIR}/swagger
	rm -rf v${SWAGGER_UI_VERSION}.tar.gz swagger-ui-${SWAGGER_UI_VERSION}
