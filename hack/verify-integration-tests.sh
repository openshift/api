#!/bin/bash

set -o nounset
set -o pipefail

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# Build yq when it's not present and not overriden for a specific file.
if [ -z "${YQ:-}" ];then
  ${TOOLS_MAKE} yq
  YQ="${TOOLS_OUTPUT}/yq"
fi

validate_suite_files() {
  FOLDER=$1

  for file in ${FOLDER}/*.yaml; do
    if [ ! -f $file ]; then
      # It's likely the bash expansion didn't find any yaml files.
      continue
    fi

    if [ $(${YQ} eval '.apiVersion' $file) != "apiextensions.k8s.io/v1" ]; then
      continue
    fi

    if [ $(${YQ} eval '.kind' $file) != "CustomResourceDefinition" ]; then
      continue
    fi

    CRD_NAME=$(echo $file | sed s:"${FOLDER}/":: )
    GROUP=$(${YQ} eval '.spec.group' $file)
    KIND=$(${YQ} eval '.spec.names.kind' $file)
    SINGULAR=$(${YQ} eval '.spec.names.singular' $file)

    FILE_BASE="stable"

    FEATURESET=$(${YQ} eval '.metadata.annotations["release.openshift.io/feature-set"]' $file)
    if [ ${FEATURESET} == "TechPreviewNoUpgrade" ]; then
      # TechPreviewNoUpgrade CRDs should start with techpreview for their test suites.
      FILE_BASE="techpreview"
    fi

    SUITE_FILE=${FOLDER}/${FILE_BASE}.${SINGULAR}.testsuite.yaml

    if [ ! -f ${SUITE_FILE} ]; then
      echo "No test suite file found for CRD ${file}"
      exit 1
    fi
  done
}

for groupVersion in ${API_GROUP_VERSIONS}; do
  echo "Validating integration tests for ${groupVersion}"
  validate_suite_files ${SCRIPT_ROOT}/${groupVersion}
done
