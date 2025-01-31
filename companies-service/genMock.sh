#!/bin/bash

go install github.com/vektra/mockery/v2@v2.42.0

mockery  --inpackage --all --dir internal/

go mod tidy