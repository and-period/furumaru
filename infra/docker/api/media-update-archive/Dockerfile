# ビルド用コンテナ
FROM golang:1.23.3-alpine as builder

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

WORKDIR /go/src/github.com/and-period/furumaru/api

RUN apk add --update --no-cache \
      git \
      make \
      tar \
      tzdata \
      wget \
      xz

ENV FFMPEG_VERSION=7.0.2
ENV FFMPEG_ARCH=arm64

RUN wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-${FFMPEG_VERSION}-${FFMPEG_ARCH}-static.tar.xz
RUN tar -Jxvf ./ffmpeg-${FFMPEG_VERSION}-${FFMPEG_ARCH}-static.tar.xz
RUN mv ./ffmpeg-${FFMPEG_VERSION}-${FFMPEG_ARCH}-static/ffmpeg /usr/bin/ffmpeg
RUN rm -rf ./ffmpeg-${FFMPEG_VERSION}-${FFMPEG_ARCH}-static.tar.xz

ADD ./api ./

RUN go mod download
RUN go build -o ./app ./hack/media-update-archive/main.go

# 実行用コンテナ
FROM alpine

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

WORKDIR /var/api

COPY --from=builder /go/src/github.com/and-period/furumaru/api/app ./app
COPY --from=builder /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /usr/share/zoneinfo/Asia/Tokyo

CMD ["./app"]
