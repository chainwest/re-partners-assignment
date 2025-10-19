#!/bin/bash

set -e

BASE_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

if ! curl -s -f "$BASE_URL/healthz" > /dev/null 2>&1; then
    echo -e "${RED}Server not running${NC}"
    exit 1
fi

run_test() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected=$5
    
    response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
        -H "Content-Type: application/json" \
        ${data:+-d "$data"})
    status=$(echo "$response" | tail -n 1)
    
    if [ "$status" -eq "$expected" ]; then
        echo -e "${GREEN}✓ $name${NC}"
    else
        echo -e "${RED}✗ $name (expected $expected, got $status)${NC}"
    fi
}

run_test "Health check" "GET" "/healthz" "" 200
run_test "Version" "GET" "/version" "" 200
run_test "Exact match" "POST" "/packs/solve" '{"sizes": [250, 500, 1000], "amount": 1000}' 200
run_test "Minimal overage" "POST" "/packs/solve" '{"sizes": [250, 500, 1000], "amount": 251}' 200
run_test "Complex case" "POST" "/packs/solve" '{"sizes": [250, 500, 1000, 2000, 5000], "amount": 12001}' 200
run_test "Negative amount" "POST" "/packs/solve" '{"sizes": [250, 500], "amount": -100}' 422
run_test "Empty sizes" "POST" "/packs/solve" '{"sizes": [], "amount": 1000}' 422
run_test "Invalid JSON" "POST" "/packs/solve" 'invalid json' 400
run_test "Wrong method" "GET" "/packs/solve" "" 405
run_test "Metrics" "GET" "/metrics" "" 200

echo -e "\n${GREEN}Tests completed${NC}"

