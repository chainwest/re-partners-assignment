#!/bin/bash

set -e

if ! command -v go &> /dev/null; then
    echo "Go not installed"
    exit 1
fi

go mod download
go mod tidy

[ ! -f .env ] && [ -f configs/.env.example ] && cp configs/.env.example .env

mkdir -p build

echo "Setup complete. Run: make help"

