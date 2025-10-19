#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

API_URL="${API_URL:-https://re-partners-api.onrender.com}"
TESTS_PASSED=0
TESTS_FAILED=0

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    ((TESTS_PASSED++))
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
    ((TESTS_FAILED++))
}

test_endpoint() {
    local name=$1
    local method=$2
    local path=$3
    local data=$4
    local expected=$5
    
    response=$(curl -s -w "\n%{http_code}" -X $method "$API_URL$path" \
        -H "Content-Type: application/json" \
        ${data:+-d "$data"})
    status=$(echo "$response" | tail -n 1)
    
    [ "$status" -eq "$expected" ] && print_success "$name" || print_error "$name"
}

test_edge_case() {
    response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/packs/solve" \
        -H "Content-Type: application/json" \
        -d '{"sizes": [250, 500, 1000, 2000, 5000], "amount": 263}')
    body=$(echo "$response" | head -n -1)
    status=$(echo "$response" | tail -n 1)
    packs=$(echo "$body" | grep -o '"packs":[0-9]*' | grep -o '[0-9]*$')
    
    [ "$status" -eq 200 ] && [ "$packs" = "2" ] && print_success "Edge case 263" || print_error "Edge case 263"
}

test_endpoint "Health" "GET" "/healthz" "" 200
test_endpoint "Version" "GET" "/version" "" 200
test_edge_case
test_endpoint "Complex case" "POST" "/packs/solve" '{"sizes": [250, 500, 1000, 2000, 5000], "amount": 12001}' 200
test_endpoint "Validation" "POST" "/packs/solve" '{"sizes": [250, 500], "amount": -1}' 422
test_endpoint "Metrics" "GET" "/metrics" "" 200
test_endpoint "Web UI" "GET" "/" "" 200

echo -e "\n${GREEN}Passed: $TESTS_PASSED${NC} | ${RED}Failed: $TESTS_FAILED${NC}"
[ $TESTS_FAILED -eq 0 ] && exit 0 || exit 1

