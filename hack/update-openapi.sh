#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Build openapi-gen when it's not present and not overriden for a specific file.
OPENAPI_GEN=${OPENAPI_GEN:-""}
if [ -z ${OPENAPI_GEN} ];then
  ${TOOLS_MAKE} openapi-gen
fi
OPENAPI_GEN=${OPENAPI_GEN:-"${TOOLS_OUTPUT}/openapi-gen"}

verify="${VERIFY:-}"
output_package="${OUTPUT_PKG:-github.com/openshift/api/openapi}"

EXT_APIS_PKG="github.com/openshift/api"

# enumerate group versions
ALL_FQ_APIS=() # e.g. k8s.io/kubernetes/pkg/apis/apps k8s.io/api/apps/v1
INT_FQ_APIS=() # e.g. k8s.io/kubernetes/pkg/apis/apps
EXT_FQ_APIS=("k8s.io/apimachinery/pkg/apis/meta/v1") # e.g. k8s.io/api/apps/v1
for GVs in ${TYPE_PACKAGE_VERSIONS}; do
  IFS=: read -r G Vs <<<"${GVs}"

  # enumerate versions
  for V in ${Vs//,/ }; do
    ALL_FQ_APIS+=("${EXT_APIS_PKG}/${G}/${V}")
    EXT_FQ_APIS+=("${EXT_APIS_PKG}/${G}/${V}")
  done
done

function codegen::join() { local IFS="$1"; shift; echo "$*"; }

echo "Generating OpenAPI definitions for ${TYPE_PACKAGE_VERSIONS} at ${output_package}"

declare -a OPENAPI_EXTRA_PACKAGES
${OPENAPI_GEN} \
         --input-dirs "$(codegen::join , "${EXT_FQ_APIS[@]}" "${OPENAPI_EXTRA_PACKAGES[@]+"${OPENAPI_EXTRA_PACKAGES[@]}"}")" \
         --input-dirs "k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/util/intstr,k8s.io/apimachinery/pkg/api/resource,k8s.io/apimachinery/pkg/version,k8s.io/api/core/v1,k8s.io/api/rbac/v1,k8s.io/api/authorization/v1" \
         --output-package "${output_package}/generated_openapi" \
         -O zz_generated.openapi \
         --go-header-file ${SCRIPT_ROOT}/hack/empty.txt \
         ${verify}

go build github.com/openshift/api/openapi/cmd/models-schema

./models-schema  | jq '.' > ../../../${output_package}/openapi.json
