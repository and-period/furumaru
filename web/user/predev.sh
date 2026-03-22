#!/bin/sh

# Nuxt ビルドキャッシュが不整合な場合にクリーンアップ
# client.precomputed.mjs は nuxt dev 起動時に vite-builder が生成するが、
# キャッシュが古いと ENOENT エラーになるため、存在チェックで検知する
NUXT_CACHE_DIR="node_modules/.cache/nuxt/.nuxt"
if [ -d "$NUXT_CACHE_DIR" ] && [ ! -f "$NUXT_CACHE_DIR/dist/server/client.precomputed.mjs" ]; then
    echo "Cleaning stale Nuxt build cache..."
    rm -rf node_modules/.cache/nuxt .nuxt .output
fi

if ! (type mkcert > /dev/null 2>&1); then
    echo "Installing mkcert..."
    brew install mkcert
fi

if [ ! -e "localhost.pem" ] && [ ! -e "localhost-key.pem" ]; then
    echo "Generating a Certificate..."
    mkcert -install
    mkcert localhost
fi
