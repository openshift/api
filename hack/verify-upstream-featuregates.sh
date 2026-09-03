#!/usr/bin/env bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

VERIFY_DIR=$(mktemp -d -t upstream-featuregates-verify-XXXXXX)

go run --mod=vendor -trimpath github.com/openshift/api/payload-command/cmd/write-upstream-featuregates --output-directory=${VERIFY_DIR} --output-filename=zz_generated.upstreamfeaturegates.go

diff -r "${VERIFY_DIR}" ./features/zz_generated.upstreamfeaturegates.go

rm -rf "${VERIFY_DIR}"
