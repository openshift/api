#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../../../k8s.io/code-generator)}

# API_GROUP_VERSION_PACKAGES is a string of <group>/<version>.
# The deepcopy gen needs a string of <group>:<version> so substitute slashes for colons.
inputArg="$(echo ${API_GROUP_VERSIONS//\//:})"

verify="${VERIFY:-}"

GOFLAGS="" bash ${CODEGEN_PKG}/generate-groups.sh "deepcopy" \
  github.com/openshift/api/generated \
  github.com/openshift/api \
  "${inputArg}" \
  --go-header-file ${SCRIPT_ROOT}/hack/empty.txt \
  ${verify}
