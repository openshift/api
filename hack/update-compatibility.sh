#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

# API_GROUP_VERSIONS is a string of <group>/<version>.
# The compatibility gen needs a comma separated list of Go packages, so prefix each entry with a comma and the
# PACKAGE_NAME, then trim the leading comma.
inputArg="$(printf ",${PACKAGE_NAME}/%s" ${API_GROUP_VERSIONS})"
inputArg="${inputArg:1}"

go run github.com/openshift/api/cmd/openshift-compatibility-gen --input-dirs "$inputArg"
