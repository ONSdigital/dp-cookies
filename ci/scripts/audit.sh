#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-cookies
  make audit
popd