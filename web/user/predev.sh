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
