# Publish Kubebuilder Tools

This is a utility to publish the kubebuilder tools archives to the kubebuilder release bucket.

The tool takes a release image and constructs 4 archives (linux/darwin x amd64/arm64) from  `installer-kube-apiserver-artifacts` and `installer-etcd-artifacts` images,
each containing the correct `etcd` and `kube-apiserver` binaries for the architecture.

The archives are then optionally published to the release buck `openshift-kubebuilder-tools` in the `openshift-gce-devel` project on GCP.

## Usage

`-pull-secret <string>`: The path to an OpenShift pull secret file containing credentials for the payload registry. New CI payloads require an entry for `quay-proxy.ci.openshift.org`; legacy payloads require `registry.ci.openshift.org`. The Quay-proxy payload and artifact images share the `openshift/ci` repository, so the tool obtains one repository-scoped bearer token before downloading them.
`-payload <string>`: The CI payload image used to create the artifacts. Supported formats are `registry.ci.openshift.org/ocp/release:<tag>`, `registry.ci.openshift.org/ocp/release-<number>:<tag>`, and `quay-proxy.ci.openshift.org/openshift/ci:rc_payload__<version>`.
`-version <string>`: The Kubernetes semantic version to represent in the archives, with the `v` prefix, e.g. `v1.37.1`. This is **not** the OpenShift payload version. Before creating archives or uploading, the tool runs the native extracted `kube-apiserver --version` and requires a matching version. Only the trailing `-dirty` suffix and build metadata are ignored; alpha/beta/rc prerelease identifiers must match.
`-output-dir <string>`: A working directory to store the archives. The binaries will be extracted here and the archives will be created here.
`-skip-upload <bool>`: Skip uploading the artifacts to GCS and updating the index file. This can be used if you are not authenticated to GCP.
`-index-file <string>`: The path to the index file that should be updated with the new archives. This is optional and will default to `./envtest-releases.yaml`.

## Archive uploads

The tool will automatically publish the artifacts to the GCS bucket once they have been created.
To set up authentication, ensure you have logged into the gcloud CLI.

```bash
gcloud auth login
```

To skip this step, add the `-skip-upload` flag to the command.

## Using the archives

To use the archives, pass the `--remote-bucket openshift-kubebuilder-tools` flag to the envtest setup command.

```makefile
ENVTEST_K8S_VERSION = 1.29.1
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
ENVTEST = go run ${PROJECT_DIR}/vendor/sigs.k8s.io/controller-runtime/tools/setup-envtest

.PHONY: test
test: ## Run only the tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path --bin-dir $(PROJECT_DIR)/bin --remote-bucket openshift-kubebuilder-tools)" ./hack/test.sh
```
