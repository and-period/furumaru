FROM amazonlinux:2023

ENV OPENAPITOOLS_VERSION 2.13.1
ENV SWAGGER_CLI_VERSION 4.0.4

RUN dnf upgrade -y
RUN dnf install -y \
  git \
  java-17-amazon-corretto \
  make \
  nodejs \
  openssh \
  tar \
  wget

RUN npm install -g yarn
RUN yarn global add @openapitools/openapi-generator-cli@v${OPENAPITOOLS_VERSION}
RUN yarn global add swagger-cli@${SWAGGER_CLI_VERSION}
