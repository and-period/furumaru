FROM debian:bookworm-slim

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

RUN apt update -&& apt install -y --no-install-recommends \
    curl \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

RUN curl --proto '=https' --tlsv1.2 -sSf https://tiup-mirrors.pingcap.com/install.sh | sh

EXPOSE 4000 2379
CMD ["tiup", "playground"]
