#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Build protoc-gen-gogo when it's not present and not overriden for a specific file.
if [ -z "${PROTOC_GEN_GOGO:-}" ]; then
  ${TOOLS_MAKE} protoc-gen-gogo
  PROTOC_GEN_GOGO="${TOOLS_OUTPUT}/protoc-gen-gogo"
fi

GENERATOR=go-to-protobuf ${SCRIPT_ROOT}/hack/update-codegen.sh
