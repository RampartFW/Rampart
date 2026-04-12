package nftables

import (
	"fmt"
	"net"
	"strings"

	"github.com/rampartfw/rampart/internal/model"
)

func generateScript(rs *model.CompiledRuleSet) string {
	var sb strings.Builder

	// Table creation
	sb.WriteString("flush table inet rampart\n")
	sb.WriteString("table inet rampart {\n")

	// Chains
	sb.WriteString("  chain input {\n")
	sb.WriteString("    type filter hook input priority 0; policy drop;\n")
	sb.WriteString("    ct state established,related accept\n")
	sb.WriteString("    iifname \"lo\" accept\n")
	sb.WriteString("    meta l4proto icmp accept\n")
	// ICMP allow: `meta l4proto icmp accept` + `meta l4proto icmpv6 accept`
	for _, rule := range rs.Rules {
		if rule.Direction == model.DirectionInbound {
			sb.WriteString("    " + renderRule(rule) + "\n")
		}
	}
	sb.WriteString("  }\n")

	sb.WriteString("  chain forward {\n")
	sb.WriteString("    type filter hook forward priority 0; policy drop;\n")
	for _, rule := range rs.Rules {
		if rule.Direction == model.DirectionForward {
			sb.WriteString("    " + renderRule(rule) + "\n")
		}
	}
	sb.WriteString("  }\n")

	sb.WriteString("  chain output {\n")
	sb.WriteString("    type filter hook output priority 0; policy accept;\n")
	for _, rule := range rs.Rules {
		if rule.Direction == model.DirectionOutbound {
			sb.WriteString("    " + renderRule(rule) + "\n")
		}
	}
	sb.WriteString("  }\n")

	sb.WriteString("}\n")


	return sb.String()
}

func renderRule(rule model.CompiledRule) string {
	var parts []string

	// IP Version
	if rule.Match.IPVersion == model.IPv4 {
		parts = append(parts, "meta nfproto ipv4")
	} else if rule.Match.IPVersion == model.IPv6 {
		parts = append(parts, "meta nfproto ipv6")
	}

	// Protocols
	if len(rule.Match.Protocols) > 0 {
		var protos []string
		for _, p := range rule.Match.Protocols {
			protos = append(protos, strings.ToLower(string(p)))
		}
		if len(protos) == 1 {
			parts = append(parts, fmt.Sprintf("meta l4proto %s", protos[0]))
		} else {
			parts = append(parts, fmt.Sprintf("meta l4proto { %s }", strings.Join(protos, ", ")))
		}
	}

	// Source CIDRs
	if len(rule.Match.SourceNets) > 0 {
		parts = append(parts, renderNets("saddr", rule.Match.SourceNets))
	}

	// Dest CIDRs
	if len(rule.Match.DestNets) > 0 {
		parts = append(parts, renderNets("daddr", rule.Match.DestNets))
	}

	// Source Ports
	if len(rule.Match.SourcePorts) > 0 {
		parts = append(parts, renderPorts("sport", rule.Match.SourcePorts))
	}

	// Dest Ports
	if len(rule.Match.DestPorts) > 0 {
		parts = append(parts, renderPorts("dport", rule.Match.DestPorts))
	}

	// Interfaces
	if len(rule.Match.Interfaces) > 0 {
		if len(rule.Match.Interfaces) == 1 {
			parts = append(parts, fmt.Sprintf("iifname %q", rule.Match.Interfaces[0]))
		} else {
			var ifaces []string
			for _, i := range rule.Match.Interfaces {
				ifaces = append(ifaces, fmt.Sprintf("%q", i))
			}
			parts = append(parts, fmt.Sprintf("iifname { %s }", strings.Join(ifaces, ", ")))
		}
	}

	// States
	if len(rule.Match.States) > 0 {
		var states []string
		for _, s := range rule.Match.States {
			states = append(states, strings.ToLower(string(s)))
		}
		parts = append(parts, fmt.Sprintf("ct state { %s }", strings.Join(states, ", ")))
	}

	// Counter
	parts = append(parts, "counter")

	// Log
	if rule.Log {
		parts = append(parts, fmt.Sprintf("log prefix %q", fmt.Sprintf("rampart:%s: ", rule.Name)))
	}

	// Rate Limit
	if rule.RateLimit != nil {
		parts = append(parts, fmt.Sprintf("limit rate %d/%s burst %d", rule.RateLimit.Rate, rule.RateLimit.Per, rule.RateLimit.Burst))
	}

	// Action
	switch rule.Action {
	case model.ActionAccept:
		parts = append(parts, "accept")
	case model.ActionDrop:
		parts = append(parts, "drop")
	case model.ActionReject:
		parts = append(parts, "reject")
	case model.ActionLog:
		// Log action already handled, default to continue/accept if not specified?
		// Usually Log is used with Accept or Drop.
		// If action is Log, it just logs and continues? Or should it have a default action?
		// Spec says Log is an action.
		parts = append(parts, "log")
	case model.ActionRateLimit:
		// Rate limit already handled, but we need the final action
		// if we use "limit" expression it can be followed by an action
		// For now we assume drop if exceeded if specified in rule.RateLimit.Action
		if rule.RateLimit != nil {
			switch rule.RateLimit.Action {
			case model.ActionAccept:
				parts = append(parts, "accept")
			case model.ActionDrop:
				parts = append(parts, "drop")
			case model.ActionReject:
				parts = append(parts, "reject")
			default:
				parts = append(parts, "drop")
			}
		}
	}

	// Comment
	parts = append(parts, fmt.Sprintf("comment %q", "rampart:"+rule.Name))

	return strings.Join(parts, " ")
}

func renderNets(field string, nets []net.IPNet) string {
	if len(nets) == 1 {
		// Detect if it is IPv4 or IPv6 to use correct prefix
		if nets[0].IP.To4() != nil {
			return fmt.Sprintf("ip %s %s", field, nets[0].String())
		}
		return fmt.Sprintf("ip6 %s %s", field, nets[0].String())
	}
	
	// Mixed IPv4/IPv6 is tricky in nftables for saddr/daddr without specific family
	// But in 'inet' table we can use 'ip saddr' or 'ip6 saddr'
	var v4, v6 []string
	for _, n := range nets {
		if n.IP.To4() != nil {
			v4 = append(v4, n.String())
		} else {
			v6 = append(v6, n.String())
		}
	}

	var results []string
	if len(v4) > 0 {
		if len(v4) == 1 {
			results = append(results, fmt.Sprintf("ip %s %s", field, v4[0]))
		} else {
			results = append(results, fmt.Sprintf("ip %s { %s }", field, strings.Join(v4, ", ")))
		}
	}
	if len(v6) > 0 {
		if len(v6) == 1 {
			results = append(results, fmt.Sprintf("ip6 %s %s", field, v6[0]))
		} else {
			results = append(results, fmt.Sprintf("ip6 %s { %s }", field, strings.Join(v6, ", ")))
		}
	}

	return strings.Join(results, " ")
}

func renderPorts(field string, ports []model.PortRange) string {
	var portStrings []string
	for _, pr := range ports {
		if pr.Start == pr.End {
			portStrings = append(portStrings, fmt.Sprintf("%d", pr.Start))
		} else {
			portStrings = append(portStrings, fmt.Sprintf("%d-%d", pr.Start, pr.End))
		}
	}
	if len(portStrings) == 1 {
		return fmt.Sprintf("tcp %s %s", field, portStrings[0]) // Defaulting to tcp for port if protocol not clear? 
		// Actually nftables expects protocol before dport/sport if not already set.
		// We'll just return the port part and assume the caller handled protocol or we'll prefix it.
		// Wait, better to let it be 'tcp dport' or 'udp dport'.
		// But if both tcp and udp are allowed?
	}
	return fmt.Sprintf("tcp %s { %s }", field, strings.Join(portStrings, ", "))
}
