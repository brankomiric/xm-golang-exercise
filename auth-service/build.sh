#!/bin/bash
VERSION=$1

rm -rf dist/auth

go build -o dist/auth \
  -ldflags="-X main.Version=$VERSION" \
  cmd/auth/main.go