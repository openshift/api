# Detect go environment to build tools into env specific directories.
# This is helpful for those running tools locally and within containers across different OS/architechture combinations.
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
TOOLS_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
OUTPUT_DIR=$(TOOLS_DIR)/_output/bin/$(GOOS)/$(GOARCH)

# The CODEGEN_VERSION is set to the commit hash of the most recent commit that touched
# the codegen folder, vendor folder or go.mod and go.sum files.
# This should contain any file that may change the output of the build for codegen.
# Append dirty if the current working tree is dirty for any of the relevant files.
CODEGEN_VERSION=$(shell git log -n 1 --pretty=format:%H -- codegen vendor go.mod go.sum; git diff HEAD --quiet codegen vendor go.mod go.sum || echo -dirty)

# Tools builds all tools.
tools: codegen controller-gen deepcopy-gen go-to-protobuf openapi-gen protoc-gen-gogo yq

clean:
	rm $(OUTPUT_DIR)/*

.PHONY: run-codegen
run-codegen: codegen
	$(OUTPUT_DIR)/codegen --base-dir $(BASE_DIR) --api-group-versions $(API_GROUP_VERSIONS) --required-feature-sets $(OPENSHIFT_REQUIRED_FEATURESETS)

# If the codegen utility is built, and the version doesn't match the latest version,
# Remove it so that it'll be rebuilt before being executed again.
.PHONY: check-codegen-version
check-codegen-version:
	@ if [ -f $(OUTPUT_DIR)/codegen ] && [ "$(CODEGEN_VERSION)" != "$(shell $(OUTPUT_DIR)/codegen --version)" ]; then \
		echo "codegen version mismatch, removing old codegen"; \
		rm $(OUTPUT_DIR)/codegen; \
	fi

#############################################
#
# BEGIN: Shortcuts to fully qualified targets
#
#############################################

.PHONY:codegen
codegen: $(OUTPUT_DIR)/codegen

.PHONY:controller-gen
controller-gen: $(OUTPUT_DIR)/controller-gen

.PHONY:deepcopy-gen
deepcopy-gen: $(OUTPUT_DIR)/deepcopy-gen

.PHONY:go-to-protobuf
go-to-protobuf: $(OUTPUT_DIR)/go-to-protobuf

.PHONY:openapi-gen
openapi-gen: $(OUTPUT_DIR)/openapi-gen

.PHONY:protoc-gen-gogo
protoc-gen-gogo: $(OUTPUT_DIR)/protoc-gen-gogo

.PHONY:yq
yq: $(OUTPUT_DIR)/yq

###########################################
#
# End: Shortcuts to fully qualified targets
#
###########################################

#####################################
#
# BEGIN: Fully qualified tool targets
#
#####################################

# Because the codegen target relies on a phony, it becomes a phony itself.
# We must check the presence of the file before rebuilding.
$(OUTPUT_DIR)/codegen: check-codegen-version
	@ if [ ! -f $(OUTPUT_DIR)/codegen ]; then \
		echo "Building codegen version ${CODEGEN_VERSION}"; \
		go build -mod=vendor -o $(OUTPUT_DIR)/codegen -ldflags="-X main.version=${CODEGEN_VERSION}" ./codegen/cmd; \
	fi

$(OUTPUT_DIR)/controller-gen:
	go build -mod=vendor -o $(OUTPUT_DIR)/controller-gen ./vendor/sigs.k8s.io/controller-tools/cmd/controller-gen

$(OUTPUT_DIR)/deepcopy-gen:
	go build -mod=vendor -o $(OUTPUT_DIR)/deepcopy-gen ./vendor/k8s.io/code-generator/cmd/deepcopy-gen

$(OUTPUT_DIR)/go-to-protobuf:
	go build -mod=vendor -o $(OUTPUT_DIR)/go-to-protobuf ./vendor/k8s.io/code-generator/cmd/go-to-protobuf

$(OUTPUT_DIR)/openapi-gen:
	go build -mod=vendor -o $(OUTPUT_DIR)/openapi-gen ./vendor/k8s.io/code-generator/cmd/openapi-gen

$(OUTPUT_DIR)/protoc-gen-gogo:
	go build -mod=vendor -o $(OUTPUT_DIR)/protoc-gen-gogo ./vendor/k8s.io/code-generator/cmd/go-to-protobuf/protoc-gen-gogo

$(OUTPUT_DIR)/yq:
	go build -mod=vendor -o $(OUTPUT_DIR)/yq ./vendor/github.com/mikefarah/yq/v4

###################################
#
# END: Fully qualified tool targets
#
###################################
