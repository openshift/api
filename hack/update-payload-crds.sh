#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

crd_globs="\
    authorization/v1/zz_generated.crd-manifests/*_config-operator_*.crd*yaml\
    config/v1/zz_generated.crd-manifests/*_config-operator_*.crd*yaml\
    machine/v1/zz_generated.crd-manifests/*.crd*yaml\
    operator/v1/zz_generated.crd-manifests//*_config-operator_*.crd*yaml\
    operator/v1alpha1/zz_generated.crd-manifests//*_config-operator_*.crd*yaml\
    quota/v1/zz_generated.crd-manifests/*_config-operator_*.crd*yaml\
    security/v1/zz_generated.crd-manifests/*_config-operator_*.crd*yaml\
    securityinternal/v1/zz_generated.crd-manifests/*_config-operator_*.crd*yaml
    "

# To allow the crd_globs to be sourced in the verify script,
# wrap the copy action to prevent it running when sourced.
if [ "$0" = "$BASH_SOURCE" ] ; then
    for f in ${crd_globs}; do
        cp "$f" "${SCRIPT_ROOT}/payload-manifests/crds/"
    done
fi
