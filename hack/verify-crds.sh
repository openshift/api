#!/bin/bash

if [ ! -f ./_output/tools/bin/yq ]; then
    mkdir -p ./_output/tools/bin
    curl -s -f -L https://github.com/mikefarah/yq/releases/download/2.4.0/yq_$(go env GOHOSTOS)_$(go env GOHOSTARCH) -o ./_output/tools/bin/yq
    chmod +x ./_output/tools/bin/yq
fi

FAILS=false
for f in `find . -name "*crd.yaml" -type f`
do
    if [[ $(./_output/tools/bin/yq r $f apiVersion) == "apiextensions.k8s.io/v1beta1" ]]; then
        if [[ $(./_output/tools/bin/yq r $f spec.validation.openAPIV3Schema.properties.metadata.description) != "null" ]]; then
            echo "Error: cannot have a metadata description in $f"
            FAILS=true
        fi

        if [[ $(./_output/tools/bin/yq r $f spec.preserveUnknownFields) != "false" ]]; then
            echo "Error: pruning not enabled (spec.preserveUnknownFields != false) in $f"
            FAILS=true
        fi
    fi
done

if [ "$FAILS" = true ] ; then
    exit 1
fi

