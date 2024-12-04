#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Build codegen-crds when it's not present and not overriden for a specific file.
if [ -z "${KAL:-}" ];then
  ${TOOLS_MAKE} kal
  KAL="${TOOLS_OUTPUT}/kal"
fi

"${KAL}" -fix ./...
