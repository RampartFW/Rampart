package engine

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/rampartfw/rampart/internal/model"
	"gopkg.in/yaml.v3"
)

var varPattern = regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)

// ParseVariablesFile parses a rampart-vars.yaml file.
func ParseVariablesFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read variables file: %w", err)
	}

	var vars model.VariablesYAML
	if err := yaml.Unmarshal(data, &vars); err != nil {
		return nil, fmt.Errorf("unmarshal variables: %w", err)
	}

	if vars.APIVersion != "rampart.dev/v1" {
		return nil, fmt.Errorf("unsupported apiVersion: %s", vars.APIVersion)
	}

	if vars.Kind != "Variables" {
		return nil, fmt.Errorf("unsupported kind: %s", vars.Kind)
	}

	return vars.Variables, nil
}

// SubstituteVars replaces ${var_name} placeholders in YAML data with values from the map.
func SubstituteVars(data []byte, vars map[string]interface{}) ([]byte, error) {
	var err error
	result := varPattern.ReplaceAllFunc(data, func(match []byte) []byte {
		name := string(match[2 : len(match)-1])
		val, ok := vars[name]
		if !ok {
			err = fmt.Errorf("unresolved variable: %s", name)
			return match
		}

		b, marshalErr := yaml.Marshal(val)
		if marshalErr != nil {
			err = fmt.Errorf("marshal variable %s: %w", name, marshalErr)
			return match
		}
		return bytes.TrimSpace(b)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
