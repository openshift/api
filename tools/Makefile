# Detect go environment to build tools into env specific directories.
# This is helpful for those running tools locally and within containers across different OS/architechture combinations.
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
TOOLS_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
OUTPUT_DIR=$(TOOLS_DIR)/_output/bin/$(GOOS)/$(GOARCH)


# Tools builds all tools.
tools: codegen controller-gen deepcopy-gen go-to-protobuf openapi-gen protoc-gen-gogo yq

clean:
	rm $(OUTPUT_DIR)/*

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

$(OUTPUT_DIR)/codegen:
	go build -mod=vendor -o $(OUTPUT_DIR)/codegen ./codegen/cmd

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
