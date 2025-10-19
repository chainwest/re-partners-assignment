#!/bin/bash

set -e

BASE_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

test_endpoint() {
    local name=$1
    local method=$2
    local path=$3
    local data=$4
    local expected=$5
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$path" -H "Content-Type: application/json" -d "$data")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$path")
    fi
    
    status=$(echo "$response" | tail -n 1)
    [ "$status" -eq "$expected" ] && echo -e "${GREEN}✓ $name${NC}" || echo -e "${RED}✗ $name${NC}"
}

test_endpoint "Health" "GET" "/healthz" "" 200
test_endpoint "Web UI" "GET" "/" "" 200
test_endpoint "Basic solve" "POST" "/packs/solve" '{"sizes":[250,500,1000],"amount":251}' 200
test_endpoint "Edge case" "POST" "/packs/solve" '{"sizes":[250,500,1000,2000,5000],"amount":263}' 200
test_endpoint "Empty sizes" "POST" "/packs/solve" '{"sizes":[],"amount":100}' 422
test_endpoint "Negative amount" "POST" "/packs/solve" '{"sizes":[250,500],"amount":-10}' 422
test_endpoint "Invalid JSON" "POST" "/packs/solve" '{invalid}' 400

echo -e "\n${GREEN}Web UI tests completed${NC}"

