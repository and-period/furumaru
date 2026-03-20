#!/bin/sh
if ! (type mkcert > /dev/null 2>&1); then
    echo "Installing mkcert..."
    brew install mkcert
fi

if [ ! -e "localhost.pem" ] && [ ! -e "localhost-key.pem" ]; then
    echo "Generating a Certificate..."
    mkcert -install
    mkcert localhost
fi

if (type node > /dev/null 2>&1) && [ -e "scripts/patch-nuxt-dev-handler.mjs" ]; then
    node ./scripts/patch-nuxt-dev-handler.mjs
fi
