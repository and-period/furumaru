.PHONY: help setup build generate lint validate

ADMIN_INPUT_DIR       := admin
ADMIN_OUTPUT_DOCS_DIR := ./../../tmp/data/swagger/admin/docs
ADMIN_OUTPUT_TS_DIR   := ./../../tmp/data/swagger/admin/types
ADMIN_WEB_TYPES_DIR   := ./../../web/admin/src/types/api

USER_INPUT_DIR       := user
USER_OUTPUT_DOCS_DIR := ./../../tmp/data/swagger/user/docs
USER_OUTPUT_TS_DIR   := ./../../tmp/data/swagger/user/types
USER_WEB_TYPES_DIR   := ./../../web/user/src/types/api

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## 環境構築
	yarn global add @openapitools/openapi-generator-cli

build: build-admin build-user ## API仕様書の生成

build-admin: ## API仕様書の生成(管理者関連)
	$(MAKE) generate TYPE=html2 INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_DOCS_DIR}
# $(MAKE) generate TYPE=typescript-axios INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_TS_DIR}
	$(MAKE) generate TYPE=typescript-fetch INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_TS_DIR}
	mkdir -p ${ADMIN_WEB_TYPES_DIR}; cp -r ${ADMIN_OUTPUT_TS_DIR}/*.ts ${ADMIN_OUTPUT_TS_DIR}/apis ${ADMIN_OUTPUT_TS_DIR}/models ${ADMIN_WEB_TYPES_DIR}

build-user: ## API仕様書の生成(購入者関連)
	$(MAKE) generate TYPE=html2 INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}
# $(MAKE) generate TYPE=typescript-axios INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_TS_DIR}
	$(MAKE) generate TYPE=typescript-fetch INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_TS_DIR}
	mkdir -p ${USER_WEB_TYPES_DIR}; cp -r ${USER_OUTPUT_TS_DIR}/*.ts ${USER_OUTPUT_TS_DIR}/apis ${USER_OUTPUT_TS_DIR}/models ${USER_WEB_TYPES_DIR}

build-html: ## API仕様書の生成(HTML形式)
	$(MAKE) generate TYPE=html2 INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=html2 INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}

build-yaml: ## API仕様書の生成(OpenAPI形式)
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}

build-typescript: ## API仕様書の生成(TypeScript型定義とクライアント生成)
# $(MAKE) generate TYPE=typescript-axios INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_TS_DIR}
# $(MAKE) generate TYPE=typescript-axios INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_TS_DIR}
	$(MAKE) generate TYPE=typescript-fetch INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_TS_DIR}
	$(MAKE) generate TYPE=typescript-fetch INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_TS_DIR}
	mkdir -p ${ADMIN_WEB_TYPES_DIR}; cp -r ${ADMIN_OUTPUT_TS_DIR}/*.ts ${ADMIN_OUTPUT_TS_DIR}/apis ${ADMIN_OUTPUT_TS_DIR}/models ${ADMIN_WEB_TYPES_DIR}
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
