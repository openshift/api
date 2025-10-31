package generation

import (
	"fmt"

	"k8s.io/gengo/v2/generator"
	"k8s.io/gengo/v2/namer"
	"k8s.io/gengo/v2/parser"
)

// Execute implements the target execution from gengo.
func Execute(p *parser.Parser, nameSystems namer.NameSystems, defaultSystem string, getTargets func(*generator.Context) []generator.Target) error {
	c, err := generator.NewContext(p, nameSystems, defaultSystem)
	if err != nil {
		return fmt.Errorf("failed making a context: %v", err)
	}

	targets := getTargets(c)
	if err := c.ExecuteTargets(targets); err != nil {
		return fmt.Errorf("failed executing generator: %v", err)
	}

	return nil
}
