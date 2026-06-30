#!/usr/bin/env bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Build protoc-gen-gogo when it's not present and not overriden for a specific file.
if [ -z "${PROTOC_GEN_GOGO:-}" ]; then
  ${TOOLS_MAKE} protoc-gen-gogo
  PROTOC_GEN_GOGO="${TOOLS_OUTPUT}/protoc-gen-gogo"
fi

# Add the protoc-gen-gogo directory to the PATH so that the generator can find it.
PROTOC_GEN_GOGO_DIR="$(dirname "${PROTOC_GEN_GOGO}")"
PATH="${PROTOC_GEN_GOGO_DIR}:${PATH}" GENERATOR=go-to-protobuf EXTRA_ARGS=--verify ${SCRIPT_ROOT}/hack/update-codegen.sh
