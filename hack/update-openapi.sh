#!/usr/bin/env bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# OUTPUT_PATH allows the verify script to generate into a different folder.
output_path="${OUTPUT_PATH:-openapi}"
output_package="${SCRIPT_ROOT}/${output_path}"

# Generate both Go and JSON OpenAPI schemas using the integrated codegen tool
GENERATOR=openapi EXTRA_ARGS=--openapi:output-package-path=${output_path} ${SCRIPT_ROOT}/hack/update-codegen.sh
