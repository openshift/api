all: build
.PHONY: all

build:
	go build github.com/openshift/api/build/...
.PHONY: build

clean:
	rm -rf _output
.PHONY: clean

update-deps:
	hack/update-deps.sh
.PHONY: update-deps

generate:
	hack/update-deepcopy.sh
	hack/update-protobuf.sh
.PHONY: generate
