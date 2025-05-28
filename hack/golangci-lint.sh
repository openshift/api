#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Build codegen-crds when it's not present and not overriden for a specific file.
if [ -z "${GOLANGCI_LINT:-}" ];then
  ${TOOLS_MAKE} golangci-lint kube-api-linter
  GOLANGCI_LINT="${TOOLS_OUTPUT}/golangci-lint"
fi

# In CI, HOME is set to / and is not writable.
# Make sure golangci-lint can create its cache.
HOME=${HOME:-"/tmp"}
if [[ ${HOME} == "/" ]]; then
  HOME="/tmp"
fi

"${GOLANGCI_LINT}" $@
