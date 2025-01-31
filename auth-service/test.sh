#!/bin/sh

set -e -u
./genMock.sh

go test ./...