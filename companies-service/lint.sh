#!/bin/sh

go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

DIR=$(cd `dirname $0` && pwd)
echo "running pre-commit hook in $DIR"
cd $DIR
golangci-lint run --config=.golangci.yml