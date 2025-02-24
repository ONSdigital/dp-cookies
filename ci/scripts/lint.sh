#!/bin/bash -eux

cwd=$(pwd)

go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5

pushd $cwd/dp-cookies
  make lint
popd
