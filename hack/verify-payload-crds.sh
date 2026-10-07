#!/usr/bin/env bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"
source "$(dirname "${BASH_SOURCE}")/update-payload-crds.sh"

files=""

# The ingress CRD is also shipped by cluster-ingress-operator. Keep this payload
# copy on v1alpha1 until that payload owner is updated after the v1 API promotion.
ingress_crd="${SCRIPT_ROOT}/payload-manifests/crds/0000_50_ingress_02_ingresses.crd.yaml"
mapfile -t ingress_versions < <(sed -n '/^  versions:/,/^[^ ]/p' "${ingress_crd}" | sed -n 's/^  - name: //p')
if [[ "${#ingress_versions[@]}" -ne 1 || "${ingress_versions[0]}" != "v1alpha1" ]]; then
    echo "Ingress payload CRD must remain v1alpha1-only until all payload owners are updated together."
    exit 1
fi

# Check there's no diff between the files in their canonical location
# and the payload-manifests location.
for f in ${crd_globs}; do
    basename=$(basename "${f}")
    files+=${basename},
    echo "Verifying diff on ${basename}"
    diff "$f" "${SCRIPT_ROOT}/payload-manifests/crds/${basename}"
done
	
files=$(echo "${files}" | tr "," "\n")

# Check that we haven't accidentally added any files that aren't tracked
# by the crd_globs into the payload CRDs folder.
for f in "${SCRIPT_ROOT}/payload-manifests/crds/"*; do
    basename=$(basename "${f}")
    if ! grep -F -q -x "${basename}" <<< "${files}"; then
        echo "Found untracked file ${basename} in payload CRD manifests.  Please add the file to crd_globs in hack/update-payload-crds.sh."
        exit 1
    fi
done
