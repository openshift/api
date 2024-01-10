#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

go run --mod=vendor -trimpath github.com/openshift/api/payload-command/cmd/write-available-featuresets --asset-output-dir=./payload-manifests/featuregates
