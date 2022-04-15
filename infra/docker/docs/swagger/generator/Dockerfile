FROM node:16.14.2-alpine

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo
ENV SHELL /bin/sh

RUN apk add --update \
  g++ \
  make \
  openjdk11-jre-headless \
  && rm -rf /var/cache/apk/*
