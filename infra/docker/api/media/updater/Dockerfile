# ビルド用コンテナ
FROM golang:1.21.0-alpine as builder

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

WORKDIR /go/src/github.com/and-period/furumaru/api

RUN apk add --update --no-cache \
      git \
      make

ADD ./api ./
RUN make build SERVICE=media/updater

# 実行用コンテナ
FROM alpine

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

WORKDIR /var/api

RUN apk add --update --no-cache \
      ca-certificates

COPY --from=builder /go/src/github.com/and-period/furumaru/api/app ./app

ADD ./api/config ./config

CMD ["./app"]
