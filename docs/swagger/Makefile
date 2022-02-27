.PHONY: build generate lint validate

USER_INPUT_DIR       := user
USER_OUTPUT_DOCS_DIR := ./../../tmp/data/swagger/user/docs
USER_OUTPUT_TS_DIR   := ./../../tmp/data/swagger/user/types
USER_WEB_TYPES_DIR   := ./../../web/user/src/types/api

build: build-user

build-user:
	$(MAKE) generate TYPE=html2 INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=typescript-axios INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_TS_DIR}
	mkdir -p ${USER_WEB_TYPES_DIR}; cp ${USER_OUTPUT_TS_DIR}/*.ts ${USER_WEB_TYPES_DIR}

build-html:
	$(MAKE) generate TYPE=html2 INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}

build-yaml:
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}

build-typescript:
	$(MAKE) generate TYPE=typescript-axios INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_TS_DIR}
	mkdir -p ${USER_WEB_TYPES_DIR}; cp ${USER_OUTPUT_TS_DIR}/*.ts ${USER_WEB_TYPES_DIR}

generate:
	openapi-generator-cli generate \
		-g ${TYPE} \
		-i ${INPUT_DIR}/openapi.yaml \
		-o ${OUTPUT_DIR}

lint:
	$(MAKE) validate INPUT_DIR=${USER_INPUT_DIR}

validate:
	openapi-generator-cli validate \
		-i ${INPUT_DIR}/openapi.yaml \
		--recommend
