#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Build openapi-gen when it's not present and not overriden for a specific file.
if [ -z "${OPENAPI_GEN:-}" ]; then
  ${TOOLS_MAKE} openapi-gen
  OPENAPI_GEN="${TOOLS_OUTPUT}/openapi-gen"
fi

verify="${VERIFY:-}"

# OUTPUT_PATH allows the verify script to generate into a different folder.
output_path="${OUTPUT_PATH:-openapi}"
output_package="${SCRIPT_ROOT}/${output_path}"

# API_GROUP_VERSIONS is a string of <group>/<version>.
# The compatibility gen needs a comma separated list of Go packages, so prefix each entry with a comma and the
# PACKAGE_NAME, then trim the leading comma.
inputArg="$(printf ",${PACKAGE_NAME}/%s" ${API_GROUP_VERSIONS})"
inputArg="${inputArg:1}"

function codegen::join() { local IFS="$1"; shift; echo "$*"; }

echo Generating OpenAPI definitions for ${API_GROUP_VERSIONS} at ${output_package}

declare -a OPENAPI_EXTRA_PACKAGES
# Clear the GOPATH and use a relative output package.
# This should make the output correct no matter whether you are in GOPATH or not.
GOPATH= ${OPENAPI_GEN} \
         --input-dirs "$(codegen::join , "${inputArg}" "${OPENAPI_EXTRA_PACKAGES[@]+"${OPENAPI_EXTRA_PACKAGES[@]}"}")" \
         --input-dirs "k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/util/intstr,k8s.io/apimachinery/pkg/api/resource,k8s.io/apimachinery/pkg/version,k8s.io/api/core/v1,k8s.io/api/rbac/v1,k8s.io/api/authorization/v1" \
         --output-package "./${output_path}/generated_openapi" \
         -O zz_generated.openapi \
         --go-header-file ${SCRIPT_ROOT}/hack/empty.txt \
         ${verify}

go build github.com/openshift/api/openapi/cmd/models-schema

./models-schema  | jq '.' > ${output_package}/openapi.json
