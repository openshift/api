#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Use PULL_BASE_REF for CI, otherwise use master unless overriden.
COMPARISON_BASE=${COMPARISON_BASE:-${PULL_BASE_SHA:-"master"}}

GENERATOR=schemacheck EXTRA_ARGS=--schemacheck:comparison-base=${COMPARISON_BASE} ${SCRIPT_ROOT}/hack/update-codegen.sh
