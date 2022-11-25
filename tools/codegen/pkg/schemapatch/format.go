package schemapatch

import (
	"bytes"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

// formatData formats the given YAML data.
// We use indentation of 2 spaces, and the yaml.v3 library as this is what
// other kube generators use.
func formatData(in []byte) ([]byte, error) {
	node := &yaml.Node{}
	if err := yaml.Unmarshal(in, node); err != nil {
		return nil, fmt.Errorf("could not unmarshal YAML: %v", err)
	}

	buf := bytes.NewBuffer(nil)
	enc := yaml.NewEncoder(buf)

	enc.SetIndent(2)

	if err := enc.Encode(node); err != nil {
		return nil, fmt.Errorf("could not encode YAML: %v", err)
	}

	return buf.Bytes(), nil
}
