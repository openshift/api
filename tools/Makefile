# Detect go environment to build tools into env specific directories.
# This is helpful for those running tools locally and within containers across different OS/architechture combinations.
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
TOOLS_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
OUTPUT_DIR=$(TOOLS_DIR)/_output/bin/$(GOOS)/$(GOARCH)


# Tools builds all tools.
tools: openapi-gen

# This acts as a shortcut target for the expanded target.
.PHONY:openapi-gen
openapi-gen: $(OUTPUT_DIR)/openapi-gen


$(OUTPUT_DIR)/openapi-gen:
	go build -mod=vendor -o $(OUTPUT_DIR)/openapi-gen ./vendor/k8s.io/code-generator/cmd/openapi-gen
