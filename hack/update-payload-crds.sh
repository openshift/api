#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

crd_globs="\
    config/v1/*_config-operator_*.crd*yaml\
    quota/v1/*.crd*yaml\
    security/v1/*.crd*yaml\
    securityinternal/v1/*.crd*yaml\
    authorization/v1/*.crd*yaml\
    operator/v1alpha1/0000_10_config-operator_01_imagecontentsourcepolicy.crd*yaml\
    operator/v1/0000_10_config-operator_*.yaml\
    config/v1alpha1/0000_10_config-operator_01_clusterimagepolicy-*.crd*yaml\
    config/v1alpha1/0000_10_config-operator_01_imagepolicy-*.crd*yaml\
    "

# To allow the crd_globs to be sourced in the verify script,
# wrap the copy action to prevent it running when sourced.
if [ "$0" = "$BASH_SOURCE" ] ; then
    for f in ${crd_globs}; do
        cp "$f" "${SCRIPT_ROOT}/payload-manifests/crds/"
    done
fi
