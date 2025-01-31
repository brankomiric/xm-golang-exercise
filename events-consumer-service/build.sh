#!/bin/bash
VERSION=$1

rm -rf dist/consumer

go build -o dist/consumer \
  -ldflags="-X main.Version=$VERSION" \
  cmd/consumer/main.go