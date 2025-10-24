#!/bin/bash

# This script verifies that all featuregates enabled in the Default featureset
# are also enabled in the OKD featureset. OKD may have additional featuregates
# (e.g., from TechPreviewNoUpgrade or DevPreviewNoUpgrade), but it must include
# all Default featuregates.

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Build codegen when it's not present and not overridden for a specific file.
if [ -z "${CODEGEN:-}" ];then
  ${TOOLS_MAKE} codegen
  CODEGEN="${TOOLS_OUTPUT}/codegen"
fi

"${CODEGEN}" verify-okd-featuregates
