#!/bin/bash

set -e

GREEN='\033[0;32m'
NC='\033[0m'

if ! docker ps | grep -q re-partners-redis; then
    docker-compose -f deployments/docker-compose.yaml up -d redis
    sleep 3
fi

go test -v ./internal/infra/redis/ 2>&1 | grep -E "^(PASS|FAIL|ok)"

echo -e "${GREEN}Redis tests completed${NC}"

