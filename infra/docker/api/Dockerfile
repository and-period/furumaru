# ビルド用コンテナ
FROM golang:1.23.0-alpine as builder

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

ARG SERVICE

WORKDIR /go/src/github.com/and-period/furumaru/api

RUN apk add --update --no-cache \
      git \
      make \
      tzdata

ADD ./api ./

RUN make build SERVICE=${SERVICE}

# 実行用コンテナ
FROM alpine

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

WORKDIR /var/api

RUN apk add --update --no-cache \
      ca-certificates \
      curl

COPY --from=builder /go/src/github.com/and-period/furumaru/api/app ./app
COPY --from=builder /go/src/github.com/and-period/furumaru/api/config ./config
COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /usr/share/zoneinfo/Asia/Tokyo

EXPOSE 8080 9090
CMD ["./app"]
