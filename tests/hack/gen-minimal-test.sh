#!/bin/bash

set -o nounset
set -o pipefail

FOLDER=$1
VERSION=$2

for file in ${FOLDER}/*.yaml; do
  if [ $(yq eval '.apiVersion' $file) != "apiextensions.k8s.io/v1" ]; then
    continue
  fi

  if [ $(yq eval '.kind' $file) != "CustomResourceDefinition" ]; then
    continue
  fi

  CRD_NAME=$(echo $file | sed s:"${FOLDER}/":: )
  GROUP=$(yq eval '.spec.group' $file)
  KIND=$(yq eval '.spec.names.kind' $file)
  SINGULAR=$(yq eval '.spec.names.singular' $file)

  SUITE_FILE=${FOLDER}/stable.${SINGULAR}.testsuite.yaml

  if [ -f ${SUITE_FILE} ]; then
    continue
  fi

  cat > ${SUITE_FILE} <<EOF
apiVersion: apiextensions.k8s.io/v1 # Hack because controller-gen complains if we don't have this
name: "[Stable] ${KIND}"
crd: ${CRD_NAME}
tests:
  onCreate:
  - name: Should be able to create a minimal ${KIND}
    initial: |
      apiVersion: ${GROUP}/${VERSION}
      kind: ${KIND}
      spec: {} # No spec is required for a ${KIND}
    expected: |
      apiVersion: ${GROUP}/${VERSION}
      kind: ${KIND}
      spec: {}
EOF

  MAKEFILE=${FOLDER}/Makefile
  if [ ! -f ${MAKEFILE} ]; then
    cat > ${MAKEFILE} <<EOF
.PHONY: test
test:
	make -C ../../tests test GINKGO_EXTRA_ARGS=--focus="${GROUP}/${VERSION}"
EOF
  fi
done
