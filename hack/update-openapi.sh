#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../../../k8s.io/code-generator)}

verify="${VERIFY:-}"
output_package="${OUTPUT_PKG:-github.com/openshift/api/openapi}"


EXT_APIS_PKG="github.com/openshift/api"

# enumerate group versions
ALL_FQ_APIS=() # e.g. k8s.io/kubernetes/pkg/apis/apps k8s.io/api/apps/v1
INT_FQ_APIS=() # e.g. k8s.io/kubernetes/pkg/apis/apps
EXT_FQ_APIS=() # e.g. k8s.io/api/apps/v1
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
go install ./${CODEGEN_PKG}/cmd/openapi-gen
declare -a OPENAPI_EXTRA_PACKAGES
${GOPATH}/bin/openapi-gen \
         --input-dirs "$(codegen::join , "${EXT_FQ_APIS[@]}" "${OPENAPI_EXTRA_PACKAGES[@]+"${OPENAPI_EXTRA_PACKAGES[@]}"}")" \
         --input-dirs "k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/version" \
         --output-package "${output_package}/generated_openapi" \
         -O zz_generated.openapi \
         --go-header-file ${SCRIPT_ROOT}/hack/empty.txt \
         ${verify}

go build github.com/openshift/api/openapi/cmd/models-schema

./models-schema  | jq '.' > ../../../${output_package}/openapi.json