package engine

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/rampartfw/rampart/internal/model"
)

type Simulator struct {
	rules []model.CompiledRule
}

func NewSimulator(rules []model.CompiledRule) *Simulator {
	return &Simulator{rules: rules}
}

// Simulate is a convenience function that creates a new simulator and runs a simulation.
func Simulate(rules []model.CompiledRule, pkt model.SimulatedPacket) model.SimulationResult {
	s := NewSimulator(rules)
	return s.Simulate(pkt)
}

func (s *Simulator) Simulate(pkt model.SimulatedPacket) model.SimulationResult {
	start := time.Now()

	for i, rule := range s.rules {
		// Skip rules for wrong direction
		if rule.Direction != pkt.Direction {
			continue
		}

		// Check schedule (is rule active now?)
		// Note: IsActive is part of Milestone 12, but we implement a basic version or skip if not available
		if rule.Schedule != nil && !IsActive(rule.Schedule, start) {
			continue
		}

		if matchesPacket(rule.Match, pkt) {
			return model.SimulationResult{
				Verdict:     rule.Action,
				MatchedRule: &s.rules[i],
				MatchPath:   buildMatchPath(rule, pkt),
				Evaluated:   i + 1,
				Duration:    time.Since(start),
			}
		}
	}

	// No rule matched -> default policy
	return model.SimulationResult{
		Verdict:   model.ActionDrop, // Default deny
		Evaluated: len(s.rules),
		Duration:  time.Since(start),
		MatchPath: "no matching rule; default policy: drop",
	}
}

func matchesPacket(match model.CompiledMatch, pkt model.SimulatedPacket) bool {
	// Protocol check
	if len(match.Protocols) > 0 {
		found := false
		for _, p := range match.Protocols {
			if p == model.ProtocolAny || p == pkt.Protocol {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Source CIDR check
	if len(match.SourceNets) > 0 {
		if !anyContains(match.SourceNets, pkt.SourceIP) {
			return false
		}
	}

	// Dest CIDR check
	if len(match.DestNets) > 0 {
		if !anyContains(match.DestNets, pkt.DestIP) {
			return false
		}
	}

	// Dest port check
	if len(match.DestPorts) > 0 {
		if !portInRanges(pkt.DestPort, match.DestPorts) {
			return false
		}
	}

	// Source port check
	if len(match.SourcePorts) > 0 {
		if !portInRanges(pkt.SourcePort, match.SourcePorts) {
			return false
		}
	}

	// Interface check
	if len(match.Interfaces) > 0 && pkt.Interface != "" {
		found := false
		for _, iface := range match.Interfaces {
			if iface == pkt.Interface {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Connection state check
	if len(match.States) > 0 && pkt.State != "" {
		found := false
		for _, state := range match.States {
			if state == pkt.State {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Negation check
	if match.Negated != nil {
		if matchesPacket(*match.Negated, pkt) {
			return false // Negated match hit -> overall miss
		}
	}

	return true
}

func anyContains(nets []net.IPNet, ip net.IP) bool {
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

func portInRanges(port uint16, ranges []model.PortRange) bool {
	for _, r := range ranges {
		if port >= r.Start && port <= r.End {
			return true
		}
	}
	return false
}

func buildMatchPath(rule model.CompiledRule, pkt model.SimulatedPacket) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("rule: %s", rule.Name))
	parts = append(parts, fmt.Sprintf("priority: %d", rule.Priority))
	parts = append(parts, fmt.Sprintf("action: %s", rule.Action))

	matchParts := []string{}
	if len(rule.Match.Protocols) > 0 {
		matchParts = append(matchParts, fmt.Sprintf("protocol: %v", rule.Match.Protocols))
	}
	if len(rule.Match.SourceNets) > 0 {
		matchParts = append(matchParts, "source IP matched")
	}
	if len(rule.Match.DestNets) > 0 {
		matchParts = append(matchParts, "dest IP matched")
	}
	if len(rule.Match.DestPorts) > 0 {
		matchParts = append(matchParts, "dest port matched")
	}

	if len(matchParts) > 0 {
		parts = append(parts, fmt.Sprintf("matches: %s", strings.Join(matchParts, ", ")))
	}

	return strings.Join(parts, "; ")
}
