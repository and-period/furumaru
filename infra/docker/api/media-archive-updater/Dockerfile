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
      wget

RUN wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz \
  && tar -xvf ./ffmpeg-release-amd64-static.tar.xz \
  && mv ./ffmpeg-release-amd64-static/ffmpeg /usr/bin/ffmpeg \
  && rn -rf ./ffmpeg-release-amd64-static

ADD ./api ./

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
