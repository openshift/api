#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

EXTRA_ARGS=--verify-only ${SCRIPT_ROOT}/hack/update-prerelease-lifecycle-gen.sh
