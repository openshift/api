//go:build tools
// +build tools

package tools

import (
	_ "github.com/gogo/protobuf/gogoproto"
	_ "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/sortkeys"
	_ "github.com/mikefarah/yq/v4"
	_ "k8s.io/code-generator"
	_ "k8s.io/code-generator/cmd/go-to-protobuf"
	_ "k8s.io/code-generator/cmd/go-to-protobuf/protoc-gen-gogo"
)
