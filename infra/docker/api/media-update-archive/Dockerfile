# ビルド用コンテナ
FROM golang:1.23.3-bookworm as builder

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

ARG SERVICE

WORKDIR /go/src/github.com/and-period/furumaru/api

RUN apt update && apt install -y --no-install-recommends \
    make \
    tzdata \
    wget \
    xz-utils \
    && rm -rf /var/lib/apt/lists/*

ENV FFMPEG_VERSION=7.0.2
ENV FFMPEG_ARCH=arm64

RUN wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-${FFMPEG_VERSION}-${FFMPEG_ARCH}-static.tar.xz \
    && tar -Jxvf ./ffmpeg-${FFMPEG_VERSION}-${FFMPEG_ARCH}-static.tar.xz \
    && mv ./ffmpeg-${FFMPEG_VERSION}-${FFMPEG_ARCH}-static/ffmpeg /usr/bin/ffmpeg \
    && rm -rf ./ffmpeg-${FFMPEG_VERSION}-${FFMPEG_ARCH}-static.tar.xz

RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    --mount=type=bind,source=./api/go.mod,target=go.mod \
    --mount=type=bind,source=./api/go.sum,target=go.sum \
    go mod download
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,source=./api/,target=. \
    make build-hack SERVICE=${SERVICE} BUILD_DIR=/dist

# 実行用コンテナ
FROM debian:bookworm-slim

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

WORKDIR /var/api

RUN apt update && apt install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /dist/app ./app
COPY --from=builder /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /usr/share/zoneinfo/Asia/Tokyo

CMD ["./app"]
