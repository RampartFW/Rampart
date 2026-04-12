package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rampartfw/rampart/internal/model"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorBold   = "\033[1m"
)

// FormatConflicts formats a list of conflicts into a string based on the requested format.
func FormatConflicts(conflicts []model.Conflict, format string) string {
	switch strings.ToLower(format) {
	case "json":
		return formatConflictsJSON(conflicts)
	default:
		return formatConflictsText(conflicts)
	}
}

func formatConflictsJSON(conflicts []model.Conflict) string {
	data, err := json.MarshalIndent(conflicts, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(data)
}

func formatConflictsText(conflicts []model.Conflict) string {
	if len(conflicts) == 0 {
		return "No conflicts detected."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%sConflict Report%s\n", colorBold, colorReset))
	sb.WriteString(strings.Repeat("=", 20) + "\n\n")

	errors := 0
	warnings := 0
	infos := 0

	for _, c := range conflicts {
		icon := ""
		color := ""
		switch c.Severity {
		case model.SeverityError:
			icon = "✗"
			color = colorRed
			errors++
		case model.SeverityWarning:
			icon = "⚠"
			color = colorYellow
			warnings++
		case model.SeverityInfo:
			icon = "ℹ"
			color = colorBlue
			infos++
		}

		sb.WriteString(fmt.Sprintf("%s%s [%s]: %s%s\n", color, icon, c.Type, c.Message, colorReset))
		sb.WriteString(fmt.Sprintf("  Rule A: %s (Priority %d, Policy %s)\n", c.RuleA.Name, c.RuleA.Priority, c.RuleA.PolicyName))
		sb.WriteString(fmt.Sprintf("    Match: %s\n", formatMatch(c.RuleA.Match)))
		sb.WriteString(fmt.Sprintf("  Rule B: %s (Priority %d, Policy %s)\n", c.RuleB.Name, c.RuleB.Priority, c.RuleB.PolicyName))
		sb.WriteString(fmt.Sprintf("    Match: %s\n", formatMatch(c.RuleB.Match)))
		
		suggestion := getSuggestion(c)
		if suggestion != "" {
			sb.WriteString(fmt.Sprintf("  %sSuggestion: %s%s\n", colorBold, suggestion, colorReset))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Summary: %d error(s), %d warning(s), %d info(s)\n", errors, warnings, infos))

	return sb.String()
}

func getSuggestion(c model.Conflict) string {
	switch c.Type {
	case model.ConflictShadow:
		return fmt.Sprintf("Consider reordering rules or narrowing the scope of %q.", c.RuleA.Name)
	case model.ConflictContradiction:
		return "Change the priority of one of the rules to resolve the contradiction."
	case model.ConflictRedundancy:
		return fmt.Sprintf("Remove rule %q as it is identical to %q.", c.RuleB.Name, c.RuleA.Name)
	case model.ConflictSubset:
		return fmt.Sprintf("Rule %q is redundant because %q already covers it with the same action.", c.RuleB.Name, c.RuleA.Name)
	case model.ConflictOverlap:
		return "Ensure the overlapping rules have different priorities to define an explicit evaluation order."
	default:
		return ""
	}
}

func formatMatch(m model.CompiledMatch) string {
	var parts []string
	if len(m.Protocols) > 0 {
		var protos []string
		for _, p := range m.Protocols {
			protos = append(protos, string(p))
		}
		parts = append(parts, fmt.Sprintf("proto %s", strings.Join(protos, ",")))
	}
	if len(m.SourceNets) > 0 {
		var nets []string
		for _, n := range m.SourceNets {
			nets = append(nets, n.String())
		}
		parts = append(parts, fmt.Sprintf("src %s", strings.Join(nets, ",")))
	}
	if len(m.DestNets) > 0 {
		var nets []string
		for _, n := range m.DestNets {
			nets = append(nets, n.String())
		}
		parts = append(parts, fmt.Sprintf("dst %s", strings.Join(nets, ",")))
	}
	if len(m.DestPorts) > 0 {
		var ports []string
		for _, p := range m.DestPorts {
			if p.Start == p.End {
				ports = append(ports, fmt.Sprintf("%d", p.Start))
			} else {
				ports = append(ports, fmt.Sprintf("%d-%d", p.Start, p.End))
			}
		}
		parts = append(parts, fmt.Sprintf("dport %s", strings.Join(ports, ",")))
	}
	if len(m.States) > 0 {
		var states []string
		for _, s := range m.States {
			states = append(states, string(s))
		}
		parts = append(parts, fmt.Sprintf("state %s", strings.Join(states, ",")))
	}
	if len(parts) == 0 {
		return "any"
	}
	return strings.Join(parts, " ")
}
