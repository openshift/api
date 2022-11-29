#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

GENERATOR=schemapatch EXTRA_ARGS=--verify ${SCRIPT_ROOT}/hack/update-codegen.sh
