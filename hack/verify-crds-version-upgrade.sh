#!/bin/bash

# This script is to verify whether api spec is kept same when switching from crd v1beta1 to v1.
# The process is:
# 1. start a kube-apiserver.
# 2. apply all the v1beta1 crds from the last commit, and dump the spec of each crd to #crdName-before file.
# 3. apply all the v1 crds from the current head, and dump the spec of each crd to #crdName-after file.
# 4. compare #crdName-before and #crdName-after for v1beta1 crd if it is switched to v1 in dev branch.

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
TMP_ROOT="${SCRIPT_ROOT}/_tmp"
KUBECTL="./_output/tools/kubebuilder/kubectl --server=http://127.0.0.1:8080"
mkdir -p "${TMP_ROOT}"

if [ ! -f ./_output/tools/bin/yq ]; then
    mkdir -p ./_output/tools/bin
    curl -s -f -L https://github.com/mikefarah/yq/releases/download/2.4.0/yq_$(go env GOHOSTOS)_$(go env GOHOSTARCH) -o ./_output/tools/bin/yq
    chmod +x ./_output/tools/bin/yq
fi

# install kube-apiserver from kubebuilder to verify apis.
if [ ! -f ./_output/tools/kubebuilder/kube-apiserver ]; then
    mkdir -p ./_output/tools/kubebuilder
	curl -f -L https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-1.19.2-$(go env GOHOSTOS)-$(go env GOHOSTARCH).tar.gz -o ./_output/tools/kubebuilder/kubebuilder.tar.gz
	tar -C ./_output/tools/kubebuilder --strip-components=2 -zvxf ./_output/tools/kubebuilder/kubebuilder.tar.gz
fi

function etcdstart {
    nohup ./_output/tools/kubebuilder/etcd --advertise-client-urls http://127.0.0.1:2379 --data-dir "${TMP_ROOT}" --listen-client-urls http://127.0.0.1:2379 1>/dev/null &
    ETCD_PID=$!
}

function apiserverstart {
    nohup ./_output/tools/kubebuilder/kube-apiserver --cert-dir $TMP_ROOT --etcd-servers http://127.0.0.1:2379 --insecure-bind-address 127.0.0.1 1>/dev/null &
    APISERVER_PID=$!
}

# Start apiserver for crd spec verification.
etcdstart
# sleep 5 second to wait for etcd start
sleep 5
apiserverstart
# sleep 5 second to wait for apiserver start
sleep 5

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
LAST_COMMIT=$(git rev-parse HEAD~1)
# Switch to master branch and apply all the existing v1beta1 crds
git checkout $LAST_COMMIT
for f in `find . -name "*crd.yaml" -type f`
do
    if [[ $(./_output/tools/bin/yq r $f apiVersion) == "apiextensions.k8s.io/v1beta1" ]]; then
        v1beta1CRDName=$(./_output/tools/bin/yq r $f metadata.name)
        v1beta1CRDNames=("${v1beta1CRDNames[*]}" $v1beta1CRDName)
        $KUBECTL apply -f $f
        $KUBECTL get crd $v1beta1CRDName -o jsonpath='{.spec}' > $TMP_ROOT/$v1beta1CRDName-before
        $KUBECTL delete -f $f
    fi
done

# Switch to current branch and apply the crd with v1 version
FALSE=false
git checkout $CURRENT_BRANCH
for f in `find . -name "*crd.yaml" -type f`
do
    if [[ $(./_output/tools/bin/yq r $f apiVersion) == "apiextensions.k8s.io/v1" ]]; then
        v1CRDName=$(./_output/tools/bin/yq r $f metadata.name)
        $KUBECTL apply -f $f || FALSE=true
        $KUBECTL get crd $v1CRDName -o jsonpath='{.spec}' > $TMP_ROOT/$v1CRDName-after
    fi
done

for v1beta1CRDName in $v1beta1CRDNames
do
    # Only compare those switch from v1beta1 to v1
    if [ -f $TMP_ROOT/$v1beta1CRDName-after ]; then
        diff -u $TMP_ROOT/$v1beta1CRDName-before $TMP_ROOT/$v1beta1CRDName-after || FALSE=true
    fi
done

function cleanup {
  if [[ -n "${ETCD_PID-}" ]]; then
    kill -9 "${ETCD_PID}" &>/dev/null || :
    wait "${ETCD_PID}" &>/dev/null || :
  fi
  if [[ -n "${APISERVER_PID-}" ]]; then
    kill -9 "${APISERVER_PID}" &>/dev/null || :
    wait "${APISERVER_PID}" &>/dev/null || :
  fi

  rm -rf $TMP_ROOT
}
trap "cleanup" EXIT SIGINT

cleanup

if [ "$FALSE" = true ] ; then
    exit 1
fi
