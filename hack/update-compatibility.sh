#!/bin/bash

apiPkgs=(
  github.com/openshift/api/apps/v1
  github.com/openshift/api/authorization/v1
  github.com/openshift/api/build/v1
  github.com/openshift/api/config/v1
  github.com/openshift/api/helm/v1beta1
  github.com/openshift/api/console/v1
  github.com/openshift/api/image/v1
  github.com/openshift/api/image/docker10
  github.com/openshift/api/image/dockerpre012
  github.com/openshift/api/imageregistry/v1
  github.com/openshift/api/kubecontrolplane/v1
  github.com/openshift/api/legacyconfig/v1
  github.com/openshift/api/network/v1
  github.com/openshift/api/oauth/v1
  github.com/openshift/api/openshiftcontrolplane/v1
  github.com/openshift/api/operator/v1
  github.com/openshift/api/operator/v1alpha1
  github.com/openshift/api/operatorcontrolplane/v1alpha1
  github.com/openshift/api/operatoringress/v1
  github.com/openshift/api/osin/v1
  github.com/openshift/api/project/v1
  github.com/openshift/api/quota/v1
  github.com/openshift/api/route/v1
  github.com/openshift/api/samples/v1
  github.com/openshift/api/security/v1
  github.com/openshift/api/securityinternal/v1
  github.com/openshift/api/servicecertsigner/v1alpha1
  github.com/openshift/api/template/v1
  github.com/openshift/api/user/v1
)
inputArg=$(printf ",%s" "${apiPkgs[@]}")
inputArg=${inputArg:1}

go run github.com/openshift/api/cmd/openshift-compatibility-gen --input-dirs "$inputArg"
