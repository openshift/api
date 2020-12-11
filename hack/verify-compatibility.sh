#!/bin/bash
#shellcheck source=lib/init.sh
source "$(dirname "${BASH_SOURCE[0]}")/lib/init.sh"

SCRIPT_ROOT=$(readlink -f "$(dirname "$(dirname "${BASH_SOURCE[0]}")")")
TMP_ROOT=$(mktemp --directory)

cleanup() {
  rm -rf "${TMP_ROOT}"
}
trap "cleanup" EXIT SIGINT

V_ROOT="${TMP_ROOT}/src/github.com/openshift/api"
mkdir -p "$V_ROOT"
cp -a --no-preserve=timestamp -r "$SCRIPT_ROOT"/* "$V_ROOT"
(
  cd "$V_ROOT" || exit
  export GOPATH="$TMP_ROOT"
  rm -Rf _output
#  bash --init-file <(echo cd $V_ROOT)
  ./hack/update-compatibility.sh > /dev/null
)

if ! diff --unified=0 --label "missing file" --label "unexpected file" \
  <(find "${V_ROOT}" -name "zz_generated.openshift_compatibility.go" -printf "%P\n" | sort) \
  <(find "${SCRIPT_ROOT}" -name "zz_generated.openshift_compatibility.go" -printf "%P\n" | sort)
then
  printf "\nopenshift_compatibility is out of date. Please run hack/update-compatibility.sh\n"
  exit 1
fi

mapfile -t GENERATED < <(find "${SCRIPT_ROOT}" -name "zz_generated.openshift_compatibility.go" -printf "%P\n" | sort)
for g in "${GENERATED[@]}" ; do
  if ! diff --unified --text "$SCRIPT_ROOT/$g" "$V_ROOT/$g" ; then
    printf "\nopenshift_compatibility is out of date. Please run hack/update-compatibility.sh\n"
    exit 1
  fi
done
