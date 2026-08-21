#!/usr/bin/env bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

rm -f ./features/zz_generated.upstreamfeaturegates.go
go run --mod=vendor -trimpath github.com/openshift/api/payload-command/cmd/write-upstream-featuregates --output-directory=./features --output-filename=zz_generated.upstreamfeaturegates.go
