#!/bin/bash
set -e -u
./genMock.sh

go get -d -v ./...

export MODE=TESTS

go test  $(go list ./... | grep -v configs | grep -v cmd/todo | grep -v /internal/database | grep -v /internal/cache) \
   -race -coverprofile cover.out -covermode atomic

perc=`go tool cover -func=cover.out | tail -n 1 | sed -Ee 's!^[^[:digit:]]+([[:digit:]]+(\.[[:digit:]]+)?)%$!\1!'`
echo "Total coverage: $perc %"
res=`echo "$perc >= 80.0" | bc`
test "$res" -eq 1 && exit 0
echo "Insufficient coverage: $perc" >&2
exit 1