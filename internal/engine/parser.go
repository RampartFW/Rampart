package engine

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/rampartfw/rampart/internal/model"
	"gopkg.in/yaml.v3"
)

// ParsePolicyFile reads and parses a PolicySet YAML file.
func ParsePolicyFile(path string) (*model.PolicySetYAML, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read policy file: %w", err)
	}

	ps, err := ParsePolicyData(data)
	if err != nil {
		return nil, err
	}
	ps.SourceFile = path
	return ps, nil
}

// ParsePolicyFileWithVars reads a PolicySet YAML file, substitutes variables, and parses it.
func ParsePolicyFileWithVars(path string, vars map[string]interface{}) (*model.PolicySetYAML, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read policy file: %w", err)
	}

	substData, err := SubstituteVars(data, vars)
	if err != nil {
		return nil, fmt.Errorf("substitute variables: %w", err)
	}

	return ParsePolicyData(substData)
}

// ParsePolicyData parses PolicySet YAML data.
func ParsePolicyData(data []byte) (*model.PolicySetYAML, error) {
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, fmt.Errorf("unmarshal policy node: %w", err)
	}

	var ps model.PolicySetYAML
	if err := node.Decode(&ps); err != nil {
		return nil, fmt.Errorf("decode policy: %w", err)
	}

	// Capture line numbers from node
	captureLineNumbers(&node, &ps)

	// Initial validation of basic fields
	if ps.APIVersion == "" {
		return nil, fmt.Errorf("apiVersion is required")
	}
	if ps.Kind == "" {
		return nil, fmt.Errorf("kind is required")
	}

	// Auto-suffix bare IPs with /32 or /128
	for i := range ps.Policies {
		for j := range ps.Policies[i].Rules {
			normalizeMatchCIDRs(&ps.Policies[i].Rules[j].Match)
		}
	}

	return &ps, nil
}

func captureLineNumbers(node *yaml.Node, ps *model.PolicySetYAML) {
	if node.Kind == yaml.DocumentNode {
		for _, content := range node.Content {
			captureLineNumbers(content, ps)
		}
		return
	}

	if node.Kind != yaml.MappingNode {
		return
	}

	// Find "policies" key
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		if key == "policies" {
			policiesNode := node.Content[i+1]
			if policiesNode.Kind == yaml.SequenceNode {
				for j, pNode := range policiesNode.Content {
					if j < len(ps.Policies) {
						ps.Policies[j].Line = pNode.Line
						captureRuleLines(pNode, &ps.Policies[j])
					}
				}
			}
		}
	}
}

func captureRuleLines(node *yaml.Node, p *model.PolicyYAML) {
	if node.Kind != yaml.MappingNode {
		return
	}

	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		if key == "rules" {
			rulesNode := node.Content[i+1]
			if rulesNode.Kind == yaml.SequenceNode {
				for j, rNode := range rulesNode.Content {
					if j < len(p.Rules) {
						p.Rules[j].Line = rNode.Line
					}
				}
			}
		}
	}
}

func normalizeMatchCIDRs(m *model.MatchYAML) {
	if m == nil {
		return
	}
	for i, cidr := range m.SourceCIDRs {
		m.SourceCIDRs[i] = normalizeCIDR(cidr)
	}
	for i, cidr := range m.DestCIDRs {
		m.DestCIDRs[i] = normalizeCIDR(cidr)
	}
	if m.Not != nil {
		normalizeMatchCIDRs(m.Not)
	}
}

func normalizeCIDR(cidr string) string {
	if strings.Contains(cidr, "/") {
		return cidr
	}
	ip := net.ParseIP(cidr)
	if ip == nil {
		return cidr // Leave invalid IP for validator to catch
	}
	if ip.To4() != nil {
		return cidr + "/32"
	}
	return cidr + "/128"
}
