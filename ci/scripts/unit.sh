#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-cookies
  make test
popd