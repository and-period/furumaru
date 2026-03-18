FROM node:22.6.0-alpine@sha256:4162c8a0f1fef9d3b003eb1fd3d8a26db46815288832aa453d829f4129d4dfd3

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo
ENV SHELL /bin/sh

RUN apk add --update \
  g++ \
  make \
  openjdk11-jre-headless \
  && rm -rf /var/cache/apk/*

RUN --mount=type=bind,source=./docs/swagger/Makefile,target=Makefile \
  make setup
