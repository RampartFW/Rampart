package iptables

import (
	"fmt"
	"strings"

	"github.com/rampartfw/rampart/internal/model"
)

// ruleToArgs translates a CompiledRule into iptables command arguments.
func ruleToArgs(rule model.CompiledRule, chain string) []string {
	args := []string{"-A", chain}

	// Protocol matching
	if len(rule.Match.Protocols) > 0 {
		// iptables supports only one protocol with -p
		// if multiple protocols are needed, we would need multiple rules
		// however, the CompiledRule usually reflects what the backend can handle.
		// For now, we'll take the first one.
		protocol := strings.ToLower(rule.Match.Protocols[0].String())
		if protocol != "any" {
			args = append(args, "-p", protocol)
		}
	}

	// Source CIDR matching
	if len(rule.Match.SourceNets) > 0 {
		var srcCIDRs []string
		for _, net := range rule.Match.SourceNets {
			srcCIDRs = append(srcCIDRs, net.String())
		}
		args = append(args, "-s", strings.Join(srcCIDRs, ","))
	}

	// Destination CIDR matching
	if len(rule.Match.DestNets) > 0 {
		var dstCIDRs []string
		for _, net := range rule.Match.DestNets {
			dstCIDRs = append(dstCIDRs, net.String())
		}
		args = append(args, "-d", strings.Join(dstCIDRs, ","))
	}

	// Destination Port matching
	if len(rule.Match.DestPorts) > 0 {
		args = append(args, "-m", "multiport", "--dports", formatPortRanges(rule.Match.DestPorts))
	}

	// Source Port matching
	if len(rule.Match.SourcePorts) > 0 {
		args = append(args, "-m", "multiport", "--sports", formatPortRanges(rule.Match.SourcePorts))
	}

	// Interface matching
	if len(rule.Match.Interfaces) > 0 {
		// iptables -i for input, -o for output.
		// For simplicity, we'll use -i if it's inbound, -o if it's outbound.
		if rule.Direction == model.DirectionInbound || rule.Direction == model.DirectionForward {
			args = append(args, "-i", rule.Match.Interfaces[0])
		} else if rule.Direction == model.DirectionOutbound {
			args = append(args, "-o", rule.Match.Interfaces[0])
		}
	}

	// Connection state matching
	if len(rule.Match.States) > 0 {
		var states []string
		for _, state := range rule.Match.States {
			states = append(states, strings.ToUpper(state.String()))
		}
		args = append(args, "-m", "state", "--state", strings.Join(states, ","))
	}

	// Rate limiting
	if rule.RateLimit != nil {
		rate := fmt.Sprintf("%d/%s", rule.RateLimit.Rate, rule.RateLimit.Per)
		args = append(args, "-m", "limit", "--limit", rate)
		if rule.RateLimit.Burst > 0 {
			args = append(args, "--limit-burst", fmt.Sprintf("%d", rule.RateLimit.Burst))
		}
	}

	// Logging
	if rule.Log {
		// TODO: Implement dual-rule generation for LOG + ACTION in iptables.
		// For now, we only handle the primary action to avoid complex rule splitting here.
	}

	// Comment for identification
	args = append(args, "-m", "comment", "--comment", fmt.Sprintf("rampart:%s", rule.Name))

	// Action (Target)
	switch rule.Action {
	case model.ActionAccept:
		args = append(args, "-j", "ACCEPT")
	case model.ActionDrop:
		args = append(args, "-j", "DROP")
	case model.ActionReject:
		args = append(args, "-j", "REJECT")
	case model.ActionLog:
		args = append(args, "-j", "LOG", "--log-prefix", fmt.Sprintf("rampart:%s: ", rule.Name))
	case model.ActionRateLimit:
		// Usually handled by the limit match above + an action.
		// If action is not specified, we'll default to ACCEPT for the limited traffic.
		args = append(args, "-j", "ACCEPT")
	}

	return args
}

func formatPortRanges(ports []model.PortRange) string {
	var ranges []string
	for _, p := range ports {
		if p.Start == p.End {
			ranges = append(ranges, fmt.Sprintf("%d", p.Start))
		} else {
			ranges = append(ranges, fmt.Sprintf("%d:%d", p.Start, p.End))
		}
	}
	return strings.Join(ranges, ",")
}
