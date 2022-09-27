# Detect go environment to build tools into env specific directories.
# This is helpful for those running tools locally and within containers across different OS/architechture combinations.
ifeq ($(GO),)
GO := go
endif
GOOS=$(shell $(GO) env GOOS)
GOARCH=$(shell $(GO) env GOARCH)
TOOLS_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
OUTPUT_DIR=$(TOOLS_DIR)/_output/bin/$(GOOS)/$(GOARCH)


# Tools builds all tools.
tools: deepcopy-gen openapi-gen


#############################################
#
# BEGIN: Shortcuts to fully qualified targets
#
#############################################

.PHONY:deepcopy-gen
deepcopy-gen: $(OUTPUT_DIR)/deepcopy-gen

.PHONY:openapi-gen
openapi-gen: $(OUTPUT_DIR)/openapi-gen

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

$(OUTPUT_DIR)/deepcopy-gen:
	$(GO) build -mod=vendor -o $(OUTPUT_DIR)/deepcopy-gen ./vendor/k8s.io/code-generator/cmd/deepcopy-gen

$(OUTPUT_DIR)/openapi-gen:
	$(GO) build -mod=vendor -o $(OUTPUT_DIR)/openapi-gen ./vendor/k8s.io/code-generator/cmd/openapi-gen

###################################
#
# END: Fully qualified tool targets
#
###################################
