package iptables

import (
	"bufio"
	"bytes"
	"fmt"
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
	_, err = exec.LookPath("ip6tables")
	if err != nil {
		// IPv6 is optional but good to have
	}
	return nil
}

func (b *IptablesBackend) CurrentState() (*model.CompiledRuleSet, error) {
	rules, err := b.getManagedRules("iptables-save")
	if err != nil {
		return nil, err
	}

	ipv6Rules, err := b.getManagedRules("ip6tables-save")
	if err == nil {
		rules = append(rules, ipv6Rules...)
	}

	return &model.CompiledRuleSet{
		Rules:      rules,
		CompiledAt: time.Now(),
	}, nil
}

func (b *IptablesBackend) getManagedRules(command string) ([]model.CompiledRule, error) {
	_, err := exec.LookPath(command)
	if err != nil {
		return nil, nil // If command doesn't exist, return no rules
	}

	cmd := exec.Command(command)
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

		// Very basic parsing for now. 
		// In a real implementation, this would be more robust.
		rule := b.parseRuleLine(line)
		rules = append(rules, rule)
	}

	return rules, nil
}

func (b *IptablesBackend) parseRuleLine(line string) model.CompiledRule {
	parts := strings.Split(line, " ")
	rule := model.CompiledRule{}

	for i := 0; i < len(parts); i++ {
		switch parts[i] {
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
		case "--comment":
			if i+1 < len(parts) {
				comment := strings.Trim(parts[i+1], "\"")
				if strings.HasPrefix(comment, "rampart:") {
					rule.Name = strings.TrimPrefix(comment, "rampart:")
				}
			}
		}
	}
	return rule
}

func (b *IptablesBackend) Apply(rs *model.CompiledRuleSet) error {
	// Chain Swap Strategy
	chains := []string{"INPUT", "FORWARD", "OUTPUT"}

	// 1. Create NEW chains
	for _, chain := range chains {
		newChain := fmt.Sprintf("%s-%s-NEW", b.chainPrefix, chain)
		b.exec("iptables", "-N", newChain)
		b.exec("ip6tables", "-N", newChain)
	}

	// 2. Populate NEW chains
	for _, rule := range rs.Rules {
		cmd := "iptables"
		if rule.Match.IPVersion == model.IPv6 {
			cmd = "ip6tables"
		}
		
		chain := b.directionToChain(rule.Direction) + "-NEW"
		args := ruleToArgs(rule, chain)
		if err := b.exec(cmd, args...); err != nil {
			b.cleanupNewChains()
			return fmt.Errorf("failed to add rule %s: %w", rule.Name, err)
		}
	}

	// 3. Swap: update jump targets in base chains
	for _, chain := range chains {
		rampartChain := fmt.Sprintf("%s-%s", b.chainPrefix, chain)
		newChain := fmt.Sprintf("%s-%s-NEW", b.chainPrefix, chain)

		// Create rampart chain if it doesn't exist
		b.exec("iptables", "-N", rampartChain)
		b.exec("ip6tables", "-N", rampartChain)

		// Ensure jumps from base chains (INPUT/FORWARD/OUTPUT) to our RAMPART chains
		// This part is tricky as it affects traffic.
		// Usually we'd want to have a permanent jump to RAMPART-INPUT at the top of INPUT.
		b.ensureJump(chain, rampartChain)

		// Now swap rules from rampartChain to newChain content
		// Efficient swap:
		// 1. Flush rampartChain
		// 2. Append all rules from newChain to rampartChain
		// Or better: 
		// 1. Jump from base to NEW
		// 2. Delete jump from base to OLD
		// 3. Rename/Swap.
		
		// Let's do the strategy from §8.1 of IMPLEMENTATION.md
		b.exec("iptables", "-I", chain, "1", "-j", newChain)
		b.exec("iptables", "-D", chain, "-j", rampartChain)
		b.exec("iptables", "-F", rampartChain)
		b.exec("iptables", "-X", rampartChain)
		b.exec("iptables", "-E", newChain, rampartChain)
		
		// Same for IPv6
		b.exec("ip6tables", "-I", chain, "1", "-j", newChain)
		b.exec("ip6tables", "-D", chain, "-j", rampartChain)
		b.exec("ip6tables", "-F", rampartChain)
		b.exec("ip6tables", "-X", rampartChain)
		b.exec("ip6tables", "-E", newChain, rampartChain)
	}

	return nil
}

func (b *IptablesBackend) directionToChain(d model.Direction) string {
	switch d {
	case model.DirectionInbound:
		return b.chainPrefix + "-INPUT"
	case model.DirectionForward:
		return b.chainPrefix + "-FORWARD"
	case model.DirectionOutbound:
		return b.chainPrefix + "-OUTPUT"
	default:
		return b.chainPrefix + "-INPUT"
	}
}

func (b *IptablesBackend) ensureJump(base, target string) {
	// Check if jump exists
	if !b.jumpExists(base, target) {
		b.exec("iptables", "-I", base, "1", "-j", target)
	}
	if !b.jumpExists6(base, target) {
		b.exec("ip6tables", "-I", base, "1", "-j", target)
	}
}

func (b *IptablesBackend) jumpExists(base, target string) bool {
	output, _ := exec.Command("iptables", "-S", base).Output()
	return strings.Contains(string(output), "-j "+target)
}

func (b *IptablesBackend) jumpExists6(base, target string) bool {
	output, _ := exec.Command("ip6tables", "-S", base).Output()
	return strings.Contains(string(output), "-j "+target)
}

func (b *IptablesBackend) cleanupNewChains() {
	chains := []string{"INPUT", "FORWARD", "OUTPUT"}
	for _, chain := range chains {
		newChain := fmt.Sprintf("%s-%s-NEW", b.chainPrefix, chain)
		b.exec("iptables", "-F", newChain)
		b.exec("iptables", "-X", newChain)
		b.exec("ip6tables", "-F", newChain)
		b.exec("ip6tables", "-X", newChain)
	}
}

func (b *IptablesBackend) exec(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	return c.Run()
}

func (b *IptablesBackend) DryRun(rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	current, err := b.CurrentState()
	if err != nil {
		return nil, err
	}
	// For now, simple plan
	return &model.ExecutionPlan{
		CurrentRuleCount: len(current.Rules),
		PlannedRuleCount: len(rs.Rules),
	}, nil
}

func (b *IptablesBackend) Rollback(snapshot *model.Snapshot) error {
	// Reconstruct state and apply
	return fmt.Errorf("rollback not fully implemented")
}

func (b *IptablesBackend) Flush() error {
	chains := []string{"INPUT", "FORWARD", "OUTPUT"}
	for _, chain := range chains {
		rampartChain := fmt.Sprintf("%s-%s", b.chainPrefix, chain)
		b.exec("iptables", "-D", chain, "-j", rampartChain)
		b.exec("iptables", "-F", rampartChain)
		b.exec("iptables", "-X", rampartChain)
		
		b.exec("ip6tables", "-D", chain, "-j", rampartChain)
		b.exec("ip6tables", "-F", rampartChain)
		b.exec("ip6tables", "-X", rampartChain)
	}
	return nil
}

func (b *IptablesBackend) Stats() (map[string]model.RuleStats, error) {
	stats := make(map[string]model.RuleStats)

	// Get stats for IPv4
	ipv4Stats, err := b.getStats("iptables")
	if err == nil {
		for k, v := range ipv4Stats {
			stats[k] = v
		}
	}

	// Get stats for IPv6
	ipv6Stats, err := b.getStats("ip6tables")
	if err == nil {
		for k, v := range ipv6Stats {
			stats[k] = v
		}
	}

	return stats, nil
}

func (b *IptablesBackend) getStats(command string) (map[string]model.RuleStats, error) {
	_, err := exec.LookPath(command)
	if err != nil {
		return nil, nil
	}

	stats := make(map[string]model.RuleStats)
	chains := []string{"INPUT", "FORWARD", "OUTPUT"}

	for _, chain := range chains {
		rampartChain := fmt.Sprintf("%s-%s", b.chainPrefix, chain)
		cmd := exec.Command(command, "-L", rampartChain, "-v", "-n", "-x")
		output, err := cmd.Output()
		if err != nil {
			continue // Chain might not exist
		}

		b.parseStatsOutput(output, stats)
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

		// iptables -L -v -n output format:
		// pkts bytes target prot opt in out source destination
		pkts := fields[0]
		bytesStr := fields[1]

		// Extract rule name from comment
		var ruleName string
		for i, field := range fields {
			if field == "/*" && i+1 < len(fields) {
				comment := fields[i+1]
				if strings.HasPrefix(comment, "rampart:") {
					ruleName = strings.TrimPrefix(comment, "rampart:")
					break
				}
			}
		}

		if ruleName != "" {
			var p, by uint64
			fmt.Sscanf(pkts, "%d", &p)
			fmt.Sscanf(bytesStr, "%d", &by)
			stats[ruleName] = model.RuleStats{
				RuleID:  ruleName,
				Packets: p,
				Bytes:   by,
			}
		}
	}
}

func (b *IptablesBackend) Close() error {
	return nil
}
