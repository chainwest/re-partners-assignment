#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

API_URL="${API_URL:-http://localhost:8080}"
MAX_RETRIES=30
RETRY_DELAY=2

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

wait_for_service() {
    local retries=0
    while [ $retries -lt $MAX_RETRIES ]; do
        if curl -s -f "$API_URL/healthz" > /dev/null 2>&1; then
            print_success "Service ready"
            return 0
        fi
        retries=$((retries + 1))
        sleep $RETRY_DELAY
    done
    print_error "Service not ready"
    exit 1
}

test_health_check() {
    local response=$(curl -s -w "\n%{http_code}" "$API_URL/healthz")
    local status=$(echo "$response" | tail -n 1)
    [ "$status" = "200" ] && print_success "Health check OK" || print_error "Health check failed"
}

test_version() {
    local response=$(curl -s -w "\n%{http_code}" "$API_URL/version")
    local status=$(echo "$response" | tail -n 1)
    [ "$status" = "200" ] && print_success "Version OK" || print_error "Version failed"
}

test_solve_simple() {
    local response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/packs/solve" \
        -H "Content-Type: application/json" \
        -d '{"sizes": [250, 500, 1000, 2000, 5000], "amount": 1}')
    local status=$(echo "$response" | tail -n 1)
    [ "$status" = "200" ] && print_success "Simple case OK" || print_error "Simple case failed"
}

test_solve_complex() {
    local response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/packs/solve" \
        -H "Content-Type: application/json" \
        -d '{"sizes": [250, 500, 1000, 2000, 5000], "amount": 12001}')
    local body=$(echo "$response" | head -n -1)
    local status=$(echo "$response" | tail -n 1)
    local packs=$(echo "$body" | grep -o '"packs":[0-9]*' | grep -o '[0-9]*$')
    [ "$status" = "200" ] && [ "$packs" = "4" ] && print_success "Complex case OK" || print_error "Complex case failed"
}

test_solve_edge_case() {
    local response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/packs/solve" \
        -H "Content-Type: application/json" \
        -d '{"sizes": [250, 500, 1000, 2000, 5000], "amount": 263}')
    local body=$(echo "$response" | head -n -1)
    local status=$(echo "$response" | tail -n 1)
    local packs=$(echo "$body" | grep -o '"packs":[0-9]*' | grep -o '[0-9]*$')
    [ "$status" = "200" ] && [ "$packs" = "2" ] && print_success "Edge case OK" || print_error "Edge case failed"
}

test_validation() {
    local response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/packs/solve" \
        -H "Content-Type: application/json" \
        -d '{"sizes": [250, 500], "amount": -1}')
    local status=$(echo "$response" | tail -n 1)
    [ "$status" = "422" ] && print_success "Validation OK" || print_error "Validation failed"
}

wait_for_service
test_health_check || true
test_version || true
test_solve_simple || true
test_solve_complex || true
test_solve_edge_case || true
test_validation || true

echo -e "\n${GREEN}Passed: $TESTS_PASSED${NC} | ${RED}Failed: $TESTS_FAILED${NC}"
[ $TESTS_FAILED -eq 0 ] && exit 0 || exit 1

