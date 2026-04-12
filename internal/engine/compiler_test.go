package engine

import (
	"testing"

	"github.com/rampartfw/rampart/internal/model"
)

func TestCompile(t *testing.T) {
	ps := &model.PolicySetYAML{
		APIVersion: "rampart.dev/v1",
		Kind:       "PolicySet",
		Metadata: model.PolicyMetadata{
			Name: "test-policy",
		},
		Defaults: &model.PolicyDefaults{
			Direction: model.DirectionInbound,
			Action:    model.ActionDrop,
		},
		Policies: []model.PolicyYAML{
			{
				Name:     "web",
				Priority: 100,
				Rules: []model.RuleYAML{
					{
						Name: "allow-http",
						Match: model.MatchYAML{
							Protocol:  "tcp",
							DestPorts: 80,
						},
						Action: model.ActionAccept,
					},
					{
						Name: "allow-https",
						Match: model.MatchYAML{
							Protocol:  "tcp",
							DestPorts: "443",
						},
						Action: model.ActionAccept,
					},
				},
			},
			{
				Name:     "ssh",
				Priority: 50,
				Rules: []model.RuleYAML{
					{
						Name: "allow-ssh",
						Match: model.MatchYAML{
							Protocol:  "tcp",
							DestPorts: []interface{}{22},
						},
						Action: model.ActionAccept,
					},
				},
			},
		},
	}

	compiled, err := Compile(ps, nil)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(compiled.Rules) != 3 {
		t.Errorf("expected 3 rules, got %d", len(compiled.Rules))
	}

	// Check sorting (priority 50 should be first)
	if compiled.Rules[0].Name != "allow-ssh" {
		t.Errorf("expected first rule to be allow-ssh (priority 50), got %s", compiled.Rules[0].Name)
	}

	// Check defaults
	if compiled.Rules[0].Direction != model.DirectionInbound {
		t.Errorf("expected direction inbound (default), got %s", compiled.Rules[0].Direction)
	}

	// Check port parsing
	if compiled.Rules[0].Match.DestPorts[0].Start != 22 {
		t.Errorf("expected port 22, got %d", compiled.Rules[0].Match.DestPorts[0].Start)
	}

	if compiled.Rules[1].Match.DestPorts[0].Start != 80 {
		t.Errorf("expected port 80, got %d", compiled.Rules[1].Match.DestPorts[0].Start)
	}

	if compiled.Rules[2].Match.DestPorts[0].Start != 443 {
		t.Errorf("expected port 443, got %d", compiled.Rules[2].Match.DestPorts[0].Start)
	}

	// Check IDs
	if compiled.Rules[0].ID == "" || compiled.Rules[1].ID == "" || compiled.Rules[2].ID == "" {
		t.Errorf("expected non-empty rule IDs")
	}

	// Check Ruleset hash
	if compiled.Hash == "" {
		t.Errorf("expected non-empty ruleset hash")
	}
}

func TestCompileWithLineNumbers(t *testing.T) {
	data := `
apiVersion: rampart.dev/v1
kind: PolicySet
metadata:
  name: test
policies:
  - name: p1
    priority: 10
    rules:
      - name: r1
        action: accept
        match:
          protocol: tcp
`
	ps, err := ParsePolicyData([]byte(data))
	if err != nil {
		t.Fatalf("ParsePolicyData failed: %v", err)
	}

	compiled, err := Compile(ps, nil)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(compiled.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(compiled.Rules))
	}

	if compiled.Rules[0].SourceLine == 0 {
		t.Errorf("expected non-zero source line")
	}
}

func TestParsePorts(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected []model.PortRange
	}{
		{80, []model.PortRange{{80, 80}}},
		{"443", []model.PortRange{{443, 443}}},
		{"8000-8080", []model.PortRange{{8000, 8080}}},
		{[]interface{}{80, "443", "8000-8080"}, []model.PortRange{{80, 80}, {443, 443}, {8000, 8080}}},
	}

	for _, tt := range tests {
		got, err := parsePorts(tt.input)
		if err != nil {
			t.Errorf("parsePorts(%v) failed: %v", tt.input, err)
			continue
		}
		if len(got) != len(tt.expected) {
			t.Errorf("parsePorts(%v) expected %d results, got %d", tt.input, len(tt.expected), len(got))
			continue
		}
		for i := range got {
			if got[i] != tt.expected[i] {
				t.Errorf("parsePorts(%v) result %d: expected %v, got %v", tt.input, i, tt.expected[i], got[i])
			}
		}
	}
}
