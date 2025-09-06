FROM amazonlinux:2023

ENV SWAGGER_UI_VERSION 5.28.1
ENV ARTIFACTS_DIR /opt

RUN dnf upgrade -y
RUN dnf install -y \
  git \
  tar \
  wget

RUN wget https://github.com/swagger-api/swagger-ui/archive/refs/tags/v${SWAGGER_UI_VERSION}.tar.gz -O /tmp/swagger-ui.tar.gz \
  && tar -zxvf /tmp/swagger-ui.tar.gz -C /tmp \
  && mkdir -p ${ARTIFACTS_DIR}/swagger \
  && mv /tmp/swagger-ui-${SWAGGER_UI_VERSION}/dist/* ${ARTIFACTS_DIR}/swagger \
  && rm -rf /tmp/swagger-ui.tar.gz /tmp/swagger-ui-${SWAGGER_UI_VERSION}

RUN sed -i 's/https:\/\/petstore.swagger.io\/v2\/swagger.json/swagger.yaml/' ${ARTIFACTS_DIR}/swagger/swagger-initializer.js
