package iptables

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

type IptablesBackend struct {
	cfg         backend.BackendConfig
	chainPrefix string
}

func init() {
	backend.Register("iptables", func(cfg backend.BackendConfig) (backend.Backend, error) {
		return &IptablesBackend{
			cfg:         cfg,
			chainPrefix: "RAMPART",
		}, nil
	})
}

func (b *IptablesBackend) Name() string {
	return "iptables"
}

func (b *IptablesBackend) Capabilities() model.BackendCapabilities {
	return model.BackendCapabilities{
		IPv4:               true,
		IPv6:               true,
		RateLimiting:       true,
		ConnectionTracking: true,
		Logging:            true,
		PerRuleCounters:    true,
		AtomicReplace:      false,
		InterfaceFiltering: true,
	}
}

func (b *IptablesBackend) Probe() error {
	_, err := exec.LookPath("iptables")
	if err != nil {
		return fmt.Errorf("iptables binary not found: %w", err)
	}
	return nil
}

func (b *IptablesBackend) CurrentState(ctx context.Context) (*model.CompiledRuleSet, error) {
	rules, err := b.getManagedRules(ctx, "iptables-save")
	if err != nil {
		return nil, err
	}

	ipv6Rules, err := b.getManagedRules(ctx, "ip6tables-save")
	if err == nil {
		rules = append(rules, ipv6Rules...)
	}

	return &model.CompiledRuleSet{
		Rules:      rules,
		CompiledAt: time.Now(),
	}, nil
}

func (b *IptablesBackend) getManagedRules(ctx context.Context, command string) ([]model.CompiledRule, error) {
	_, err := exec.LookPath(command)
	if err != nil {
		return nil, nil // If command doesn't exist, return no rules
	}

	cmd := exec.CommandContext(ctx, command)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run %s: %w", command, err)
	}

	return b.parseIptablesSave(output)
}

func (b *IptablesBackend) parseIptablesSave(data []byte) ([]model.CompiledRule, error) {
	var rules []model.CompiledRule
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "-A ") {
			continue
		}

		if !strings.Contains(line, "rampart:") {
			continue
		}

		rule := b.parseRuleLine(line)
		rules = append(rules, rule)
	}

	return rules, nil
}

func (b *IptablesBackend) parseRuleLine(line string) model.CompiledRule {
	parts := strings.Split(line, " ")
	rule := model.CompiledRule{}

	for i := 0; i < len(parts); i++ {
		p := parts[i]
		switch p {
		case "-A":
			if i+1 < len(parts) {
				chain := parts[i+1]
				if strings.Contains(chain, "INPUT") {
					rule.Direction = model.DirectionInbound
				} else if strings.Contains(chain, "FORWARD") {
					rule.Direction = model.DirectionForward
				} else if strings.Contains(chain, "OUTPUT") {
					rule.Direction = model.DirectionOutbound
				}
			}
		case "-p":
			if i+1 < len(parts) {
				rule.Match.Protocols = append(rule.Match.Protocols, model.ProtocolFromString(parts[i+1]))
			}
		case "-s":
			if i+1 < len(parts) {
				_, n, _ := net.ParseCIDR(parts[i+1])
				if n == nil {
					ip := net.ParseIP(parts[i+1])
					if ip != nil {
						mask := net.CIDRMask(32, 32)
						if ip.To4() == nil {
							mask = net.CIDRMask(128, 128)
						}
						n = &net.IPNet{IP: ip, Mask: mask}
					}
				}
				if n != nil {
					rule.Match.SourceNets = append(rule.Match.SourceNets, *n)
				}
			}
		case "-d":
			if i+1 < len(parts) {
				_, n, _ := net.ParseCIDR(parts[i+1])
				if n == nil {
					ip := net.ParseIP(parts[i+1])
					if ip != nil {
						mask := net.CIDRMask(32, 32)
						if ip.To4() == nil {
							mask = net.CIDRMask(128, 128)
						}
						n = &net.IPNet{IP: ip, Mask: mask}
					}
				}
				if n != nil {
					rule.Match.DestNets = append(rule.Match.DestNets, *n)
				}
			}
		case "--dport":
			if i+1 < len(parts) {
				p := parts[i+1]
				if strings.Contains(p, ":") {
					var start, end uint16
					fmt.Sscanf(p, "%d:%d", &start, &end)
					rule.Match.DestPorts = append(rule.Match.DestPorts, model.PortRange{Start: start, End: end})
				} else {
					var port uint16
					fmt.Sscanf(p, "%d", &port)
					rule.Match.DestPorts = append(rule.Match.DestPorts, model.PortRange{Start: port, End: port})
				}
			}
		case "--sport":
			if i+1 < len(parts) {
				p := parts[i+1]
				if strings.Contains(p, ":") {
					var start, end uint16
					fmt.Sscanf(p, "%d:%d", &start, &end)
					rule.Match.SourcePorts = append(rule.Match.SourcePorts, model.PortRange{Start: start, End: end})
				} else {
					var port uint16
					fmt.Sscanf(p, "%d", &port)
					rule.Match.SourcePorts = append(rule.Match.SourcePorts, model.PortRange{Start: port, End: port})
				}
			}
		case "-m":
			if i+1 < len(parts) && parts[i+1] == "comment" {
				if i+3 < len(parts) && parts[i+2] == "--comment" {
					comment := strings.Trim(parts[i+3], "\"")
					if strings.HasPrefix(comment, "rampart:") {
						rule.Name = strings.TrimPrefix(comment, "rampart:")
					}
				}
			}
		case "-j":
			if i+1 < len(parts) {
				rule.Action = model.ActionFromString(parts[i+1])
			}
		}
	}
	return rule
}

func (b *IptablesBackend) Apply(ctx context.Context, rs *model.CompiledRuleSet) error {
	// Chain Swap Strategy
	chains := []string{"INPUT", "FORWARD", "OUTPUT"}

	for _, ipt := range []string{"iptables", "ip6tables"} {
		if _, err := exec.LookPath(ipt); err != nil {
			continue
		}

		for _, chain := range chains {
			rampartChain := fmt.Sprintf("%s-%s", b.chainPrefix, chain)
			newChain := fmt.Sprintf("%s-%s-NEW", b.chainPrefix, chain)

			// 1. Create NEW chain
			_ = b.exec(ctx, ipt, "-N", newChain)
			_ = b.exec(ctx, ipt, "-F", newChain)

			// 2. Populate NEW chain
			for _, rule := range rs.Rules {
				if !b.shouldApplyRule(rule, chain, ipt == "ip6tables") {
					continue
				}
				args := ruleToArgs(rule, newChain)
				if err := b.exec(ctx, ipt, append([]string{"-A"}, args...)...); err != nil {
					return fmt.Errorf("failed to add rule %s to %s: %w", rule.Name, newChain, err)
				}
			}

			// 3. Ensure rampart chain exists
			_ = b.exec(ctx, ipt, "-N", rampartChain)

			// 4. Atomic-ish swap using temporary jump
			// a. Insert jump to NEW at top of base chain
			_ = b.exec(ctx, ipt, "-I", chain, "1", "-j", newChain)
			// b. Delete jump to OLD from base chain
			_ = b.exec(ctx, ipt, "-D", chain, "-j", rampartChain)
			// c. Flush and delete OLD
			_ = b.exec(ctx, ipt, "-F", rampartChain)
			_ = b.exec(ctx, ipt, "-X", rampartChain)
			// d. Rename NEW to OLD
			if err := b.exec(ctx, ipt, "-E", newChain, rampartChain); err != nil {
				return fmt.Errorf("failed to rename chain %s to %s: %w", newChain, rampartChain, err)
			}
			// e. Re-add jump to rampartChain (now contains new rules)
			_ = b.exec(ctx, ipt, "-I", chain, "1", "-j", rampartChain)
			// f. Remove temporary jump to newChain (which no longer exists)
			_ = b.exec(ctx, ipt, "-D", chain, "-j", newChain)
		}
	}

	return nil
}

func (b *IptablesBackend) shouldApplyRule(rule model.CompiledRule, chain string, isIPv6 bool) bool {
	// Direction check
	switch chain {
	case "INPUT":
		if rule.Direction != model.DirectionInbound { return false }
	case "FORWARD":
		if rule.Direction != model.DirectionForward { return false }
	case "OUTPUT":
		if rule.Direction != model.DirectionOutbound { return false }
	}

	// IP Version check
	hasIPv4 := false
	hasIPv6 := false
	
	if len(rule.Match.SourceNets) == 0 && len(rule.Match.DestNets) == 0 {
		return true
	}

	for _, n := range append(rule.Match.SourceNets, rule.Match.DestNets...) {
		if n.IP.To4() != nil {
			hasIPv4 = true
		} else {
			hasIPv6 = true
		}
	}

	if isIPv6 { return hasIPv6 }
	return hasIPv4
}

func (b *IptablesBackend) exec(ctx context.Context, cmd string, args ...string) error {
	c := exec.CommandContext(ctx, cmd, args...)
	return c.Run()
}

func (b *IptablesBackend) DryRun(ctx context.Context, rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	current, err := b.CurrentState(ctx)
	if err != nil {
		return nil, err
	}
	return &model.ExecutionPlan{
		CurrentRuleCount: len(current.Rules),
		PlannedRuleCount: len(rs.Rules),
	}, nil
}

func (b *IptablesBackend) Rollback(ctx context.Context, snapshot *model.Snapshot) error {
	return fmt.Errorf("rollback not fully implemented")
}

func (b *IptablesBackend) Flush(ctx context.Context) error {
	chains := []string{"INPUT", "FORWARD", "OUTPUT"}
	for _, ipt := range []string{"iptables", "ip6tables"} {
		if _, err := exec.LookPath(ipt); err != nil {
			continue
		}
		for _, chain := range chains {
			rampartChain := fmt.Sprintf("%s-%s", b.chainPrefix, chain)
			_ = b.exec(ctx, ipt, "-D", chain, "-j", rampartChain)
			_ = b.exec(ctx, ipt, "-F", rampartChain)
			_ = b.exec(ctx, ipt, "-X", rampartChain)
		}
	}
	return nil
}

func (b *IptablesBackend) Stats(ctx context.Context) (map[string]model.RuleStats, error) {
	stats := make(map[string]model.RuleStats)
	for _, ipt := range []string{"iptables", "ip6tables"} {
		if _, err := exec.LookPath(ipt); err != nil {
			continue
		}
		for _, chain := range []string{"INPUT", "FORWARD", "OUTPUT"} {
			rampartChain := fmt.Sprintf("%s-%s", b.chainPrefix, chain)
			cmd := exec.CommandContext(ctx, ipt, "-L", rampartChain, "-v", "-n", "-x")
			output, err := cmd.Output()
			if err == nil {
				b.parseStatsOutput(output, stats)
			}
		}
	}
	return stats, nil
}

func (b *IptablesBackend) parseStatsOutput(output []byte, stats map[string]model.RuleStats) {
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.Contains(line, "rampart:") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		var pkts, bytesVal uint64
		fmt.Sscanf(fields[0], "%d", &pkts)
		fmt.Sscanf(fields[1], "%d", &bytesVal)

		for _, field := range fields {
			if strings.HasPrefix(field, "rampart:") {
				ruleName := strings.TrimPrefix(field, "rampart:")
				s := stats[ruleName]
				s.RuleID = ruleName
				s.Packets += pkts
				s.Bytes += bytesVal
				stats[ruleName] = s
				break
			}
		}
	}
}

func (b *IptablesBackend) Close() error {
	return nil
}
