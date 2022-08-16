#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

VERIFY=--verify-only ${SCRIPT_ROOT}/hack/update-deepcopy.sh
