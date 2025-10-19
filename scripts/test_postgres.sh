#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

DB_HOST=${TEST_DB_HOST:-localhost}
DB_PORT=${TEST_DB_PORT:-5432}
DB_USER=${TEST_DB_USER:-postgres}
DB_PASSWORD=${TEST_DB_PASSWORD:-postgres}
DB_NAME=${TEST_DB_NAME:-re_partners_test}

if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${RED}PostgreSQL unavailable${NC}"
    exit 1
fi

PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME" > /dev/null 2>&1
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME" > /dev/null 2>&1

PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f deployments/migrations/001_create_pack_sets.up.sql > /dev/null 2>&1
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f deployments/migrations/002_create_calculations.up.sql > /dev/null 2>&1

TEST_DB_HOST=$DB_HOST TEST_DB_PORT=$DB_PORT TEST_DB_USER=$DB_USER TEST_DB_PASSWORD=$DB_PASSWORD TEST_DB_NAME=$DB_NAME \
go test -v ./internal/infra/postgres/... 2>&1 | grep -E "PASS|FAIL"

PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME" > /dev/null 2>&1

echo -e "${GREEN}PostgreSQL tests completed${NC}"

