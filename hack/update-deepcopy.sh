#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Build deepcopy-gen when it's not present and not overriden for a specific file.
if [ -z "${DEEPCOPY_GEN:-}" ];then
  ${TOOLS_MAKE} deepcopy-gen
  DEEPCOPY_GEN="${TOOLS_OUTPUT}/deepcopy-gen"
fi

# API_GROUP_VERSION_PACKAGES is a string of <group>/<version>.
# The deepcopy gen needs a string of <group>:<version> so substitute slashes for colons.
inputArg="$(echo ${API_GROUP_VERSIONS//\//:})"

verify="${VERIFY:-}"

# API_GROUP_VERSION_PACKAGES is a string of <group>/<version>.
# The deepcopy gen needs a comma-separated list of ./<group>/<version> so print in that format and remove the leading comma.
inputArg="$(printf ",./%s" ${API_GROUP_VERSIONS})"
inputArg="${inputArg:1}"

echo Generating Deepcopy for ${API_GROUP_VERSIONS}

# Explicitly clear the GOPATH so that the output is relative.
# This ensures the command works both inside and outside of GOPATH.
GOPATH= "${DEEPCOPY_GEN}" \
  -O zz_generated.deepcopy \
  --trim-path-prefix "${PACKAGE_NAME}" \
  --go-header-file "${SCRIPT_ROOT}/hack/empty.txt" \
  --input-dirs "${inputArg}" \
  ${verify}
