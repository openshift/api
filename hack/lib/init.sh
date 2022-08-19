#!/bin/bash

# This script is meant to be the entrypoint for OpenShift Bash scripts to import all of the support
# libraries at once in order to make Bash script preambles as minimal as possible. This script recur-
# sively `source`s *.sh files in this directory tree. As such, no files should be `source`ed outside
# of this script to ensure that we do not attempt to overwrite read-only variables.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(cd $(dirname "${BASH_SOURCE[0]}")/../.. && pwd)
PACKAGE_NAME="github.com/openshift/api"

TOOLS_MAKE="make -C ${SCRIPT_ROOT}/tools"
TOOLS_OUTPUT="${SCRIPT_ROOT}/tools/_output/bin/$(go env GOOS)/$(go env GOARCH)"

# find_api_group_versions uses regex to look for any folder that looks like it would be an API group version.
# eg. foo/v1, bar/v1beta1 or baz/v1alpha1. It then outputs a list of <api>/<version> for those API group versions.
# API group versions are required for the following generators:
# - compatibility
# - deepcopy
# - openapi
# - swagger
find_api_group_versions() {
  # Use separate regexes because the `|` operator doesn't work consistently on different versions of find.
  # Use sed to trim the SCRIPT_ROOT from the output to create the <api>/<version> strings.
  find "${SCRIPT_ROOT}" -type d \( -regex "${SCRIPT_ROOT}/[a-z]*/v[0-9]" -o -regex "${SCRIPT_ROOT}/[a-z]*/v[0-9]alpha[0-9]" -o -regex "${SCRIPT_ROOT}/[a-z]*/v[0-9]beta[0-9]" \) | sed "s|^${SCRIPT_ROOT}/||" | sort
}

# Find the API group versions when not already set.
# Include non-standard groupversions from the image API as well.
API_GROUP_VERSIONS=${API_GROUP_VERSIONS:-"$(find_api_group_versions) image/docker10 image/dockerpre012"}

# API Packages is used by the protobuf generator.
# Protobuf generation is explicitly opt-in for packages so these must be listed here.
API_PACKAGES="\
github.com/openshift/api/apps/v1,\
github.com/openshift/api/authorization/v1,\
github.com/openshift/api/build/v1,\
github.com/openshift/api/image/v1,\
github.com/openshift/api/cloudnetwork/v1,\
github.com/openshift/api/network/v1,\
github.com/openshift/api/networkoperator/v1,\
github.com/openshift/api/oauth/v1,\
github.com/openshift/api/project/v1,\
github.com/openshift/api/quota/v1,\
github.com/openshift/api/route/v1,\
github.com/openshift/api/samples/v1,\
github.com/openshift/api/security/v1,\
github.com/openshift/api/template/v1,\
github.com/openshift/api/user/v1\
"
