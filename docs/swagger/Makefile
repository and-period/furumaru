.PHONY: build generate lint validate

ADMIN_INPUT_DIR       := admin
ADMIN_OUTPUT_DOCS_DIR := ./../../tmp/data/swagger/admin/docs
ADMIN_OUTPUT_TS_DIR   := ./../../tmp/data/swagger/admin/types
ADMIN_WEB_TYPES_DIR   := ./../../web/admin/src/types/api

USER_INPUT_DIR       := user
USER_OUTPUT_DOCS_DIR := ./../../tmp/data/swagger/user/docs
USER_OUTPUT_TS_DIR   := ./../../tmp/data/swagger/user/types
USER_WEB_TYPES_DIR   := ./../../web/user/src/types/api

build: build-admin build-user

build-admin:
	$(MAKE) generate TYPE=html2 INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=typescript-axios INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_TS_DIR}
	mkdir -p ${ADMIN_WEB_TYPES_DIR}; cp ${ADMIN_OUTPUT_TS_DIR}/*.ts ${ADMIN_WEB_TYPES_DIR}

build-user:
	$(MAKE) generate TYPE=html2 INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=typescript-axios INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_TS_DIR}
	mkdir -p ${USER_WEB_TYPES_DIR}; cp ${USER_OUTPUT_TS_DIR}/*.ts ${USER_WEB_TYPES_DIR}

build-html:
	$(MAKE) generate TYPE=html2 INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=html2 INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}

build-yaml:
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_DOCS_DIR}
	$(MAKE) generate TYPE=openapi-yaml INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_DOCS_DIR}

build-typescript:
	$(MAKE) generate TYPE=typescript-axios INPUT_DIR=${ADMIN_INPUT_DIR} OUTPUT_DIR=${ADMIN_OUTPUT_TS_DIR}
	$(MAKE) generate TYPE=typescript-axios INPUT_DIR=${USER_INPUT_DIR} OUTPUT_DIR=${USER_OUTPUT_TS_DIR}
	mkdir -p ${ADMIN_WEB_TYPES_DIR}; cp ${ADMIN_OUTPUT_TS_DIR}/*.ts ${ADMIN_WEB_TYPES_DIR}
	mkdir -p ${USER_WEB_TYPES_DIR}; cp ${USER_OUTPUT_TS_DIR}/*.ts ${USER_WEB_TYPES_DIR}

generate:
	openapi-generator-cli generate \
		-g ${TYPE} \
		-i ${INPUT_DIR}/openapi.yaml \
		-o ${OUTPUT_DIR}

lint:
	$(MAKE) validate INPUT_DIR=${ADMIN_INPUT_DIR}
	$(MAKE) validate INPUT_DIR=${USER_INPUT_DIR}

validate:
	openapi-generator-cli validate \
		-i ${INPUT_DIR}/openapi.yaml \
		--recommend
