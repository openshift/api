#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Build prerelease-lifecycle-gen when it's not present and not overriden for a specific file.
if [ -z "${PRERELEASE_LIFECYCLE_GEN:-}" ];then
  ${TOOLS_MAKE} prerelease-lifecycle-gen
  PRERELEASE_LIFECYCLE_GEN="${TOOLS_OUTPUT}/prerelease-lifecycle-gen"
fi

"${PRERELEASE_LIFECYCLE_GEN}" --logtostderr -v 1 -h $(dirname "${BASH_SOURCE}")/boilerplate.go.txt --input-dirs ${API_PACKAGES} ${EXTRA_ARGS:-}
