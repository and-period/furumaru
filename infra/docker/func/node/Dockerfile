FROM node:20.11.1-alpine

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

RUN apk add --update \
  g++ \
  jq \
  make \
  python3 \
  zip \
  && ln -fs python3 /usr/bin/python
