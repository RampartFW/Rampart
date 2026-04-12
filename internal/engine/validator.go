package engine

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/rampartfw/rampart/internal/model"
)

// ValidatePolicySet performs programmatic validation of a PolicySetYAML.
func ValidatePolicySet(ps *model.PolicySetYAML) error {
	var errs []error

	if ps.APIVersion != "rampart.dev/v1" {
		errs = append(errs, fmt.Errorf("unsupported apiVersion: %s (must be rampart.dev/v1)", ps.APIVersion))
	}
	if ps.Kind != "PolicySet" {
		errs = append(errs, fmt.Errorf("unsupported kind: %s (must be PolicySet)", ps.Kind))
	}
	if ps.Metadata.Name == "" {
		errs = append(errs, fmt.Errorf("metadata.name is required"))
	}

	policyNames := make(map[string]bool)
	for i, p := range ps.Policies {
		if p.Name == "" {
			errs = append(errs, fmt.Errorf("policies[%d].name is required", i))
		} else if policyNames[p.Name] {
			errs = append(errs, fmt.Errorf("duplicate policy name: %s", p.Name))
		}
		policyNames[p.Name] = true

		if p.Priority < 0 || p.Priority > 999 {
			errs = append(errs, fmt.Errorf("policy %s: priority must be 0-999 (got %d)", p.Name, p.Priority))
		}

		ruleNames := make(map[string]bool)
		for j, r := range p.Rules {
			if r.Name == "" {
				errs = append(errs, fmt.Errorf("policy %s, rules[%d].name is required", p.Name, j))
			} else if ruleNames[r.Name] {
				errs = append(errs, fmt.Errorf("policy %s: duplicate rule name: %s", p.Name, r.Name))
			}
			ruleNames[r.Name] = true

			if err := validateRuleYAML(&r); err != nil {
				errs = append(errs, fmt.Errorf("policy %s, rule %s: %w", p.Name, r.Name, err))
			}
		}
	}

	return errors.Join(errs...)
}

func validateRuleYAML(r *model.RuleYAML) error {
	var errs []error

	if err := validateMatchYAML(&r.Match); err != nil {
		errs = append(errs, err)
	}

	if r.Action == "" {
		errs = append(errs, fmt.Errorf("action is required"))
	} else if !isValidAction(r.Action) {
		errs = append(errs, fmt.Errorf("invalid action: %s", r.Action))
	}

	if r.RateLimit != nil {
		if r.RateLimit.Rate <= 0 {
			errs = append(errs, fmt.Errorf("rateLimit.rate must be > 0"))
		}
		if r.RateLimit.Burst < r.RateLimit.Rate {
			errs = append(errs, fmt.Errorf("rateLimit.burst must be >= rate"))
		}
	}

	if r.Schedule != nil {
		if err := validateScheduleYAML(r.Schedule); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func validateMatchYAML(m *model.MatchYAML) error {
	var errs []error

	if m.Protocol != nil {
		if err := validateProtocols(m.Protocol); err != nil {
			errs = append(errs, err)
		}
	}

	for _, cidr := range m.SourceCIDRs {
		if !isValidCIDR(cidr) {
			errs = append(errs, fmt.Errorf("invalid source CIDR: %s", cidr))
		}
	}

	for _, cidr := range m.DestCIDRs {
		if !isValidCIDR(cidr) {
			errs = append(errs, fmt.Errorf("invalid destination CIDR: %s", cidr))
		}
	}

	if m.SourcePorts != nil {
		if err := validatePorts(m.SourcePorts); err != nil {
			errs = append(errs, fmt.Errorf("invalid sourcePorts: %w", err))
		}
	}

	if m.DestPorts != nil {
		if err := validatePorts(m.DestPorts); err != nil {
			errs = append(errs, fmt.Errorf("invalid destPorts: %w", err))
		}
	}

	if m.Not != nil {
		if err := validateMatchYAML(m.Not); err != nil {
			errs = append(errs, fmt.Errorf("not match: %w", err))
		}
	}

	return errors.Join(errs...)
}

func validateProtocols(p interface{}) error {
	switch v := p.(type) {
	case string:
		if !isValidProtocol(v) {
			return fmt.Errorf("invalid protocol: %s", v)
		}
	case []interface{}:
		for _, item := range v {
			s, ok := item.(string)
			if !ok {
				return fmt.Errorf("protocol must be string")
			}
			if !isValidProtocol(s) {
				return fmt.Errorf("invalid protocol: %s", s)
			}
		}
	default:
		return fmt.Errorf("protocol must be string or []string")
	}
	return nil
}

func isValidProtocol(p string) bool {
	p = strings.ToLower(p)
	switch p {
	case "tcp", "udp", "icmp", "icmpv6", "any":
		return true
	}
	return false
}

func isValidAction(a model.Action) bool {
	switch a {
	case model.ActionAccept, model.ActionDrop, model.ActionReject, model.ActionLog, model.ActionRateLimit:
		return true
	}
	return false
}

func isValidCIDR(cidr string) bool {
	if !strings.Contains(cidr, "/") {
		return net.ParseIP(cidr) != nil
	}
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

func validatePorts(p interface{}) error {
	switch v := p.(type) {
	case int:
		if v < 1 || v > 65535 {
			return fmt.Errorf("port must be 1-65535")
		}
	case []interface{}:
		for _, item := range v {
			if err := validatePorts(item); err != nil {
				return err
			}
		}
	case string:
		// Check for range
		if strings.Contains(v, "-") {
			parts := strings.Split(v, "-")
			if len(parts) != 2 {
				return fmt.Errorf("invalid port range: %s", v)
			}
			var start, end int
			if _, err := fmt.Sscanf(parts[0], "%d", &start); err != nil {
				return fmt.Errorf("invalid start port: %s", parts[0])
			}
			if _, err := fmt.Sscanf(parts[1], "%d", &end); err != nil {
				return fmt.Errorf("invalid end port: %s", parts[1])
			}
			if start < 1 || start > 65535 || end < 1 || end > 65535 || start > end {
				return fmt.Errorf("invalid port range values: %s", v)
			}
		} else {
			var port int
			if _, err := fmt.Sscanf(v, "%d", &port); err != nil {
				return fmt.Errorf("invalid port: %s", v)
			}
			if port < 1 || port > 65535 {
				return fmt.Errorf("port must be 1-65535")
			}
		}
	default:
		return fmt.Errorf("invalid port type: %T", p)
	}
	return nil
}

func validateScheduleYAML(s *model.ScheduleYAML) error {
	var from, until time.Time
	var err error

	if s.ActiveFrom != "" {
		from, err = time.Parse(time.RFC3339, s.ActiveFrom)
		if err != nil {
			return fmt.Errorf("invalid activeFrom format: %w", err)
		}
	}

	if s.ActiveUntil != "" {
		until, err = time.Parse(time.RFC3339, s.ActiveUntil)
		if err != nil {
			return fmt.Errorf("invalid activeUntil format: %w", err)
		}
	}

	if !from.IsZero() && !until.IsZero() && !from.Before(until) {
		return fmt.Errorf("activeFrom must be before activeUntil")
	}

	return nil
}
