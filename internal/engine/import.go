package engine

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rampartfw/rampart/internal/model"
)

func ImportIptablesSave(data []byte) (*model.PolicySetYAML, error) {
	ps := &model.PolicySetYAML{
		APIVersion: "rampartfw.com/v1",
		Kind:       "PolicySet",
		Metadata: model.PolicyMetadata{
			Name:        "imported-iptables",
			Description: "Imported from iptables-save",
		},
		Policies: []model.PolicyYAML{},
	}

	policyMap := make(map[string]*model.PolicyYAML)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "-A ") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		chain := parts[1]
		policy, ok := policyMap[chain]
		if !ok {
			direction := model.DirectionInbound
			switch chain {
			case "INPUT":
				direction = model.DirectionInbound
			case "FORWARD":
				direction = model.DirectionForward
			case "OUTPUT":
				direction = model.DirectionOutbound
			}

			policy = &model.PolicyYAML{
				Name:      strings.ToLower(chain),
				Priority:  500,
				Direction: direction,
				Rules:     []model.RuleYAML{},
			}
			policyMap[chain] = policy
			ps.Policies = append(ps.Policies, *policy)
		}

		rule := model.RuleYAML{
			Name:  fmt.Sprintf("rule-%d", len(policy.Rules)+1),
			Match: model.MatchYAML{},
		}

		for i := 2; i < len(parts); i++ {
			switch parts[i] {
			case "-p":
				rule.Match.Protocol = parts[i+1]
				i++
			case "-s":
				rule.Match.SourceCIDRs = append(rule.Match.SourceCIDRs, parts[i+1])
				i++
			case "-d":
				rule.Match.DestCIDRs = append(rule.Match.DestCIDRs, parts[i+1])
				i++
			case "--dport":
				rule.Match.DestPorts = parts[i+1]
				i++
			case "--sport":
				rule.Match.SourcePorts = parts[i+1]
				i++
			case "-i":
				rule.Match.Interfaces = append(rule.Match.Interfaces, parts[i+1])
				i++
			case "-j":
				rule.Action = model.Action(strings.ToLower(parts[i+1]))
				i++
			case "-m":
				if parts[i+1] == "comment" && parts[i+2] == "--comment" {
					rule.Description = strings.Trim(parts[i+3], "\"")
					i += 3
				}
			}
		}

		// Update the policy in the slice (policy points to the one in policyMap)
		// but since we appended it to the slice, we need to update it in the slice too.
		// Actually it's easier to just rebuild the ps.Policies at the end.
		policy.Rules = append(policy.Rules, rule)
	}

	ps.Policies = []model.PolicyYAML{}
	for _, p := range policyMap {
		ps.Policies = append(ps.Policies, *p)
	}

	return ps, nil
}

type nftJSON struct {
	Nftables []json.RawMessage `json:"nftables"`
}

type nftRule struct {
	Rule struct {
		Family  string            `json:"family"`
		Table   string            `json:"table"`
		Chain   string            `json:"chain"`
		Expr    []json.RawMessage `json:"expr"`
		Comment string            `json:"comment"`
	} `json:"rule"`
}

func ImportNftables(data []byte) (*model.PolicySetYAML, error) {
	var nJSON nftJSON
	if err := json.Unmarshal(data, &nJSON); err != nil {
		return nil, fmt.Errorf("failed to parse nftables JSON: %w", err)
	}

	ps := &model.PolicySetYAML{
		APIVersion: "rampartfw.com/v1",
		Kind:       "PolicySet",
		Metadata: model.PolicyMetadata{
			Name:        "imported-nftables",
			Description: "Imported from nft list",
		},
		Policies: []model.PolicyYAML{},
	}

	policyMap := make(map[string]*model.PolicyYAML)

	for _, raw := range nJSON.Nftables {
		var r nftRule
		if err := json.Unmarshal(raw, &r); err != nil || r.Rule.Chain == "" {
			continue
		}

		policyName := fmt.Sprintf("%s-%s", r.Rule.Table, r.Rule.Chain)
		policy, ok := policyMap[policyName]
		if !ok {
			policy = &model.PolicyYAML{
				Name:     policyName,
				Priority: 500,
				Rules:    []model.RuleYAML{},
			}
			policyMap[policyName] = policy
		}

		rule := model.RuleYAML{
			Name:        fmt.Sprintf("rule-%d", len(policy.Rules)+1),
			Description: r.Rule.Comment,
			Match:       model.MatchYAML{},
		}

		var exprs []map[string]interface{}
		if err := json.Unmarshal(json.RawMessage(fmt.Sprintf("[%s]", strings.Join(func() []string {
			var s []string
			for _, e := range r.Rule.Expr {
				s = append(s, string(e))
			}
			return s
		}(), ","))), &exprs); err != nil {
			continue
		}

		for _, expr := range exprs {
			if match, ok := expr["match"].(map[string]interface{}); ok {
				left := match["left"].(map[string]interface{})
				right := match["right"]
				op := match["op"].(string)

				if op == "==" {
					payload, ok := left["payload"].(map[string]interface{})
					if ok {
						protocol := payload["protocol"].(string)
						field := payload["field"].(string)

						switch field {
						case "saddr":
							rule.Match.SourceCIDRs = append(rule.Match.SourceCIDRs, fmt.Sprint(right))
						case "daddr":
							rule.Match.DestCIDRs = append(rule.Match.DestCIDRs, fmt.Sprint(right))
						case "dport":
							rule.Match.DestPorts = right
						case "sport":
							rule.Match.SourcePorts = right
						case "protocol":
							rule.Match.Protocol = right
						}

						if protocol != "ip" && protocol != "ip6" {
							rule.Match.Protocol = protocol
						}
					}
				}
			} else if verdict, ok := expr["accept"]; ok {
				_ = verdict
				rule.Action = model.ActionAccept
			} else if verdict, ok := expr["drop"]; ok {
				_ = verdict
				rule.Action = model.ActionDrop
			}
			// Add more expr handlers as needed
		}

		policy.Rules = append(policy.Rules, rule)
	}

	for _, p := range policyMap {
		ps.Policies = append(ps.Policies, *p)
	}

	return ps, nil
}

