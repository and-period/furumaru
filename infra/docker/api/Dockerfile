# ビルド用コンテナ
FROM golang:1.25.3-alpine as builder

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

ARG SERVICE

WORKDIR /go/src/github.com/and-period/furumaru/api

RUN apk add --update --no-cache \
      git \
      make \
      tzdata

# 依存関係のダウンロードを別レイヤーにして最適化
COPY ./api/go.mod ./api/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    go mod download

# ビルド処理
COPY ./api ./
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    make build SERVICE=${SERVICE} BUILD_DIR=/dist

# 実行用コンテナ
FROM alpine

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

WORKDIR /var/api

RUN apk add --update --no-cache \
      ca-certificates \
      curl

COPY ./api/config ./config

COPY --from=builder /dist/app ./app
COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /usr/share/zoneinfo/Asia/Tokyo

EXPOSE 8080 9090
CMD ["./app"]
