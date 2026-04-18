package engine

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/rampartfw/rampart/internal/model"
	"gopkg.in/yaml.v3"
)

var (
	varPattern     = regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)
	dynamicPattern = regexp.MustCompile(`\$\{([a-z0-9]+):([^}]+)\}`)
)

// DynamicResolver defines the interface for sources that can resolve variables at runtime.
type DynamicResolver interface {
	Name() string
	Resolve(ctx context.Context, query string) ([]string, error)
}

var resolvers = make(map[string]DynamicResolver)

// RegisterResolver adds a new dynamic resolver to the engine.
func RegisterResolver(r DynamicResolver) {
	resolvers[r.Name()] = r
}

// ResolveDynamicVars processes strings like ${k8s:app=frontend} into actual values.
func ResolveDynamicVars(ctx context.Context, input string) ([]string, error) {
	match := dynamicPattern.FindStringSubmatch(input)
	if match == nil {
		return []string{input}, nil
	}

	source := match[1]
	query := match[2]

	resolver, ok := resolvers[source]
	if !ok {
		return nil, fmt.Errorf("unknown dynamic source: %s", source)
	}

	return resolver.Resolve(ctx, query)
}

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

	if vars.APIVersion != "rampartfw.com/v1" {
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

