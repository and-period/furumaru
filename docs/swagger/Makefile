.PHONY: help setup clean build generate

OPENAPITOOLS_VERSION := 2.13.4
SWAGGER_CLI_VERSION  := 4.0.4
SWAGGER_UI_VERSION   := 5.12.0

ADMIN_V1_INPUT         := gen/admin/v1/swagger.yaml
ADMIN_V1_OUTPUT_TS_DIR := ./../../tmp/data/swagger/admin/v1/types

USER_V1_INPUT         := gen/user/v1/swagger.yaml
USER_V1_OUTPUT_TS_DIR := ./../../tmp/data/swagger/user/v1/types

USER_FACILITY_INPUT         := gen/user/facility/swagger.yaml
USER_FACILITY_OUTPUT_TS_DIR := ./../../tmp/data/swagger/user/facility/types

ADMIN_WEB_V1_TYPES_DIR  := ./../../web/admin/src/types/api/v1
USER_WEB_V1_TYPES_DIR   := ./../../web/user/src/types/api/v1
LIFF_V1_TYPES_DIR       := ./../../web/liff/src/types/api/v1
LIFF_FACILITY_TYPES_DIR := ./../../web/liff/src/types/api/facility

ARTIFACTS_DIR := artifacts

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## 環境構築
	yarn global add \
		@openapitools/openapi-generator-cli@v${OPENAPITOOLS_VERSION} \
		swagger-cli@v${SWAGGER_CLI_VERSION}

clean: ## 生成されたファイルを削除
	rm -rf ${ADMIN_WEB_V1_TYPES_DIR}
	rm -rf ${USER_WEB_V1_TYPES_DIR}
	rm -rf ${LIFF_V1_TYPES_DIR}
	rm -rf ${LIFF_FACILITY_TYPES_DIR}

build: clean build-admin-v1 build-user-v1 build-user-facility ## API仕様書の生成

build-admin-v1: ## API仕様書の生成(管理者関連)
	$(MAKE) generate TYPE=typescript-axios INPUT=${ADMIN_V1_INPUT} OUTPUT_DIR=${ADMIN_V1_OUTPUT_TS_DIR}
	mkdir -p ${ADMIN_WEB_V1_TYPES_DIR}; cp -r ${ADMIN_V1_OUTPUT_TS_DIR}/*.ts ${ADMIN_WEB_V1_TYPES_DIR}

build-user-v1: ## API仕様書の生成(購入者関連)
	$(MAKE) generate TYPE=typescript-fetch INPUT=${USER_V1_INPUT} OUTPUT_DIR=${USER_V1_OUTPUT_TS_DIR}
	mkdir -p ${USER_WEB_V1_TYPES_DIR}; cp -r ${USER_V1_OUTPUT_TS_DIR}/*.ts ${USER_V1_OUTPUT_TS_DIR}/apis ${USER_V1_OUTPUT_TS_DIR}/models ${USER_WEB_V1_TYPES_DIR}
	mkdir -p ${LIFF_V1_TYPES_DIR}; cp -r ${USER_V1_OUTPUT_TS_DIR}/*.ts ${USER_V1_OUTPUT_TS_DIR}/apis ${USER_V1_OUTPUT_TS_DIR}/models ${LIFF_V1_TYPES_DIR}

build-user-facility: ## API仕様書の生成(購入者施設関連)
	$(MAKE) generate TYPE=typescript-fetch INPUT=${USER_FACILITY_INPUT} OUTPUT_DIR=${USER_FACILITY_OUTPUT_TS_DIR}
	mkdir -p ${LIFF_FACILITY_TYPES_DIR}; cp -r ${USER_FACILITY_OUTPUT_TS_DIR}/*.ts ${USER_FACILITY_OUTPUT_TS_DIR}/apis ${USER_FACILITY_OUTPUT_TS_DIR}/models ${LIFF_FACILITY_TYPES_DIR}

build-docs: ## ホスティング環境へデプロイ用の成果物を生成
	$(MAKE) download-swagger-ui OUTPUT_DIR=${ARTIFACTS_DIR}
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${INPUT} OUTPUT_DIR=${ARTIFACTS_DIR}
	swagger-cli bundle ${ARTIFACTS_DIR}/openapi/openapi.yaml  --outfile ${ARTIFACTS_DIR}/swagger/openapi.yaml --type yaml --format 3
	sed -i 's/https:\/\/petstore.swagger.io\/v2\/swagger.json/openapi.yaml/' ${ARTIFACTS_DIR}/swagger/swagger-initializer.js

generate: ## API仕様書の生成
	openapi-generator-cli generate \
		-g ${TYPE} \
		-i ${INPUT} \
		-o ${OUTPUT_DIR}

download-swagger-ui: ## ホスティング時に使用するSwaggerUIのダウンロード
	mkdir -p artifacts/swagger
	wget https://github.com/swagger-api/swagger-ui/archive/refs/tags/v${SWAGGER_UI_VERSION}.tar.gz
	tar -zxvf v${SWAGGER_UI_VERSION}.tar.gz swagger-ui-${SWAGGER_UI_VERSION}/dist/
	cp swagger-ui-${SWAGGER_UI_VERSION}/dist/* ${OUTPUT_DIR}/swagger
	rm -rf v${SWAGGER_UI_VERSION}.tar.gz swagger-ui-${SWAGGER_UI_VERSION}
