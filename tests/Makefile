# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.25
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
ifeq ($(GO),)
GO := go
endif
ENVTEST = $(GO) run ${PROJECT_DIR}/vendor/sigs.k8s.io/controller-runtime/tools/setup-envtest

.PHONY: integration
integration: verify-vendor test

.PHONY: test
test: ## Run only the tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path --bin-dir $(PROJECT_DIR)/bin)" ./hack/test.sh

.PHONY: vendor
vendor:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify

.PHONY: verify-vendor
verify-vendor: vendor
	git diff --exit-code vendor
