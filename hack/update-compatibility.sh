#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

GENERATOR=compatibility ${SCRIPT_ROOT}/hack/update-codegen.sh
