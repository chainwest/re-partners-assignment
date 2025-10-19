#!/bin/bash

set -e

go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out | tail -1
go tool cover -html=coverage.out -o coverage.html

echo "Tests completed. Coverage: coverage.html"

