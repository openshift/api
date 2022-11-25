#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

GENERATOR=schemapatch ${SCRIPT_ROOT}/hack/update-codegen.sh
