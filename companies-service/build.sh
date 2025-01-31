#!/bin/bash
VERSION=$1

rm -rf dist/companies

go build -o dist/companies \
  -ldflags="-X main.Version=$VERSION" \
  cmd/companies/main.go