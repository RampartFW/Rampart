package engine

import (
	"fmt"
	"net"
	"strings"

	"github.com/rampartfw/rampart/internal/model"
)

// ValidatePolicySet performs deep semantic validation on a PolicySet.
func ValidatePolicySet(ps *model.PolicySetYAML) error {
	if ps.Metadata.Name == "" {
		return fmt.Errorf("policy set name is required")
	}

	names := make(map[string]bool)
	for i, p := range ps.Policies {
		if p.Name == "" {
			return fmt.Errorf("policy [%d] name is required", i)
		}
		if names[p.Name] {
			return fmt.Errorf("duplicate policy name: %s", p.Name)
		}
		names[p.Name] = true

		if err := validatePolicy(&p); err != nil {
			return fmt.Errorf("policy %q: %w", p.Name, err)
		}
	}

	return nil
}

func validatePolicy(p *model.PolicyYAML) error {
	ruleNames := make(map[string]bool)
	for i, r := range p.Rules {
		if r.Name == "" {
			return fmt.Errorf("rule [%d] name is required", i)
		}
		if ruleNames[r.Name] {
			return fmt.Errorf("duplicate rule name: %s", r.Name)
		}
		ruleNames[r.Name] = true

		if err := validateRule(&r); err != nil {
			return fmt.Errorf("rule %q: %w", r.Name, err)
		}
	}
	return nil
}

func validateRule(r *model.RuleYAML) error {
	// Action validation
	action := strings.ToLower(string(r.Action))
	if action != "accept" && action != "drop" && action != "reject" && action != "log" {
		return fmt.Errorf("invalid action: %s", r.Action)
	}

	// Direction validation
	dir := strings.ToLower(r.Direction)
	if dir != "" && dir != "inbound" && dir != "outbound" && dir != "forward" {
		return fmt.Errorf("invalid direction: %s", r.Direction)
	}

	// Match validation
	if err := validateMatch(&r.Match); err != nil {
		return err
	}

	return nil
}

func validateMatch(m *model.MatchYAML) error {
	// CIDR validation
	for _, cidr := range m.SourceCIDRs {
		if _, _, err := net.ParseCIDR(cidr); err != nil {
			return fmt.Errorf("invalid source CIDR: %s", cidr)
		}
	}
	for _, cidr := range m.DestCIDRs {
		if _, _, err := net.ParseCIDR(cidr); err != nil {
			return fmt.Errorf("invalid dest CIDR: %s", cidr)
		}
	}

	// Protocol validation
	if m.Protocol != nil {
		switch p := m.Protocol.(type) {
		case string:
			proto := strings.ToLower(p)
			if proto != "" && proto != "tcp" && proto != "udp" && proto != "icmp" && proto != "icmpv6" && proto != "any" {
				return fmt.Errorf("invalid protocol: %s", p)
			}
		case []string:
			for _, protoStr := range p {
				proto := strings.ToLower(protoStr)
				if proto != "tcp" && proto != "udp" && proto != "icmp" && proto != "icmpv6" {
					return fmt.Errorf("invalid protocol in list: %s", protoStr)
				}
			}
		default:
			return fmt.Errorf("protocol must be a string or a list of strings")
		}
	}

	return nil
}
