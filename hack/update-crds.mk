GO ?=go
GOHOSTOS ?=$(shell $(GO) env GOHOSTOS)
GOHOSTARCH ?=$(shell $(GO) env GOHOSTARCH)

PERMANENT_TMP :=_output
PERMANENT_TMP_GOPATH :=$(PERMANENT_TMP)/tools

CONTROLLER_GEN_VERSION ?=v0.2.1
CONTROLLER_GEN_TEMP ?=$(PERMANENT_TMP_GOPATH)/src/sigs.k8s.io/controller-tools
controller_gen_gopath =$(shell realpath -m $(CONTROLLER_GEN_TEMP)/../..)
CONTROLLER_GEN ?=$(CONTROLLER_GEN_TEMP)/controller-gen

ensure-controller-gen:
ifeq "" "$(wildcard $(CONTROLLER_GEN))"
	$(info Installing controller-gen into "$(CONTROLLER_GEN)")
	mkdir -p '$(CONTROLLER_GEN_TEMP)'
	git clone -b '$(CONTROLLER_GEN_VERSION)' --single-branch --depth=1 https://github.com/kubernetes-sigs/controller-tools.git '$(CONTROLLER_GEN_TEMP)'
	@echo '$(CONTROLLER_GEN_TEMP)/../..'
	cd '$(CONTROLLER_GEN_TEMP)' && export GO111MODULE=on GOPATH='$(controller_gen_gopath)' && $(GO) mod vendor 2>/dev/null && $(GO) build -mod=vendor ./cmd/controller-gen
else
	$(info Using existing controller-gen from "$(CONTROLLER_GEN)")
endif
.PHONY: ensure-controller-gen

clean-controller-gen:
	if [ -d '$(controller_gen_gopath)/pkg/mod' ]; then chmod +w -R '$(controller_gen_gopath)/pkg/mod'; fi
	$(RM) -r '$(CONTROLLER_GEN_TEMP)' '$(controller_gen_gopath)/pkg/mod'
	@mkdir -p '$(CONTROLLER_GEN_TEMP)'  # to make sure we can do the next step and to avoid using '/*' wildcard on the line above which could go crazy on wrong substitution
	if [ -d '$(CONTROLLER_GEN_TEMP)' ]; then rmdir --ignore-fail-on-non-empty -p '$(CONTROLLER_GEN_TEMP)'; fi
	@mkdir -p '$(controller_gen_gopath)/pkg/mod'  # to make sure we can do the next step and to avoid using '/*' wildcard on the line above which could go crazy on wrong substitution
	if [ -d '$(controller_gen_gopath)/pkg/mod' ]; then rmdir --ignore-fail-on-non-empty -p '$(controller_gen_gopath)/pkg/mod'; fi
.PHONY: clean-controller-gen

clean: clean-controller-gen

YQ ?=$(PERMANENT_TMP_GOPATH)/bin/yq
yq_dir :=$(dir $(YQ))


ensure-yq:
ifeq "" "$(wildcard $(YQ))"
	$(info Installing yq into '$(YQ)')
	mkdir -p '$(yq_dir)'
	curl -s -f -L https://github.com/mikefarah/yq/releases/download/2.4.0/yq_$(GOHOSTOS)_$(GOHOSTARCH) -o '$(YQ)'
	chmod +x '$(YQ)';
else
	$(info Using existing yq from "$(YQ)")
endif
.PHONY: ensure-yq

clean-yq:
	$(RM) '$(YQ)'
	if [ -d '$(yq_dir)' ]; then rmdir --ignore-fail-on-non-empty -p '$(yq_dir)'; fi
.PHONY: clean-yq

clean: clean-yq

crd_patches =$(subst $(CRD_SCHEMA_GEN_MANIFESTS),$(CRD_SCHEMA_GEN_OUTPUT),$(wildcard $(CRD_SCHEMA_GEN_MANIFESTS)/*.crd.yaml-merge-patch))

# $1 - crd file
# $2 - patch file
define patch-crd
	cp -n $(2) '$(CRD_SCHEMA_GEN_OUTPUT)/' || true
	$(YQ) m -i '$(1)' '$(2)'

endef

define update-crds
$(eval $(call crd-schema-gen,$(1),$(2)))
endef

empty :=
define crd-schema-gen
CRD_SCHEMA_GEN_OUTPUT =$(2)
CRD_SCHEMA_GEN_MANIFESTS =$(2)
update-codegen-crds-$(1): ensure-controller-gen ensure-yq
	'$(CONTROLLER_GEN)' \
		schemapatch:manifests="$(2)" \
		paths="$(2)" \
		output:dir="$$(CRD_SCHEMA_GEN_OUTPUT)"
	$(foreach p,$(crd_patches),$(call patch-crd,$(basename $(p)).yaml,$(p)))
.PHONY: update-codegen-crds-$(1)

verify-codegen-crds-$(1): CRD_SCHEMA_GEN_OUTPUT := $$(shell mktemp -d)
verify-codegen-crds-$(1): update-codegen-crds-$(1)
	$(foreach p,$(wildcard $$(2)/*.crd.yaml),$(call diff-crd,$(p),$(subst $$(2),$(CRD_SCHEMA_GEN_OUTPUT),$(p))))
.PHONY: verify-codegen-crds-$(1)
endef

# $1 - manifest (actual) crd
# $2 - temp crd
define diff-crd
	diff -Naup $(1) $(2)

endef
