#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Build codegen-crds when it's not present and not overriden for a specific file.
if [ -z "${CODEGEN:-}" ];then
  ${TOOLS_MAKE} codegen
  CODEGEN="${TOOLS_OUTPUT}/codegen"
fi

if [ -z "${OPENSHIFT_REQUIRED_FEATURESETS:-}" ];then
  echo "Generating CRDs..."
  "${CODEGEN}" schemapatch --base-dir "${SCRIPT_ROOT}" -v 1
else
  echo "Generating CRDs for ${OPENSHIFT_REQUIRED_FEATURESETS} FeatureSet..."
  "${CODEGEN}" schemapatch --base-dir "${SCRIPT_ROOT}" -v 1 --required-feature-sets ${OPENSHIFT_REQUIRED_FEATURESETS}
fi
