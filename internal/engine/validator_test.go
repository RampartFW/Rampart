package engine

import (
	"testing"

	"github.com/rampartfw/rampart/internal/model"
)

func TestValidatePolicySet_Valid(t *testing.T) {
	ps := &model.PolicySetYAML{
		APIVersion: "rampart.dev/v1",
		Kind:       "PolicySet",
		Metadata: model.PolicyMetadata{
			Name: "test-policy",
		},
		Policies: []model.PolicyYAML{
			{
				Name:     "web-access",
				Priority: 500,
				Rules: []model.RuleYAML{
					{
						Name: "allow-http",
						Match: model.MatchYAML{
							Protocol:  "tcp",
							DestPorts: []interface{}{80, 443},
						},
						Action: model.ActionAccept,
					},
				},
			},
		},
	}

	if err := ValidatePolicySet(ps); err != nil {
		t.Fatalf("ValidatePolicySet failed: %v", err)
	}
}

func TestValidatePolicySet_Invalid(t *testing.T) {
	tests := []struct {
		name string
		ps   *model.PolicySetYAML
	}{
		{
			name: "wrong-apiVersion",
			ps: &model.PolicySetYAML{
				APIVersion: "wrong.dev/v1",
				Kind:       "PolicySet",
				Metadata: model.PolicyMetadata{
					Name: "test",
				},
			},
		},
		{
			name: "wrong-kind",
			ps: &model.PolicySetYAML{
				APIVersion: "rampart.dev/v1",
				Kind:       "WrongKind",
				Metadata: model.PolicyMetadata{
					Name: "test",
				},
			},
		},
		{
			name: "missing-metadata-name",
			ps: &model.PolicySetYAML{
				APIVersion: "rampart.dev/v1",
				Kind:       "PolicySet",
			},
		},
		{
			name: "duplicate-policy-name",
			ps: &model.PolicySetYAML{
				APIVersion: "rampart.dev/v1",
				Kind:       "PolicySet",
				Metadata: model.PolicyMetadata{Name: "test"},
				Policies: []model.PolicyYAML{
					{Name: "p1"},
					{Name: "p1"},
				},
			},
		},
		{
			name: "invalid-priority",
			ps: &model.PolicySetYAML{
				APIVersion: "rampart.dev/v1",
				Kind:       "PolicySet",
				Metadata: model.PolicyMetadata{Name: "test"},
				Policies: []model.PolicyYAML{
					{Name: "p1", Priority: 1000},
				},
			},
		},
		{
			name: "duplicate-rule-name",
			ps: &model.PolicySetYAML{
				APIVersion: "rampart.dev/v1",
				Kind:       "PolicySet",
				Metadata: model.PolicyMetadata{Name: "test"},
				Policies: []model.PolicyYAML{
					{
						Name: "p1",
						Rules: []model.RuleYAML{
							{Name: "r1", Action: "accept"},
							{Name: "r1", Action: "drop"},
						},
					},
				},
			},
		},
		{
			name: "invalid-action",
			ps: &model.PolicySetYAML{
				APIVersion: "rampart.dev/v1",
				Kind:       "PolicySet",
				Metadata: model.PolicyMetadata{Name: "test"},
				Policies: []model.PolicyYAML{
					{
						Name: "p1",
						Rules: []model.RuleYAML{
							{Name: "r1", Action: "invalid"},
						},
					},
				},
			},
		},
		{
			name: "invalid-cidr",
			ps: &model.PolicySetYAML{
				APIVersion: "rampart.dev/v1",
				Kind:       "PolicySet",
				Metadata: model.PolicyMetadata{Name: "test"},
				Policies: []model.PolicyYAML{
					{
						Name: "p1",
						Rules: []model.RuleYAML{
							{
								Name: "r1",
								Action: "accept",
								Match: model.MatchYAML{
									SourceCIDRs: []string{"invalid-cidr"},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePolicySet(tt.ps); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestValidatePorts(t *testing.T) {
	tests := []struct {
		name    string
		ports   interface{}
		wantErr bool
	}{
		{"single-int", 80, false},
		{"out-of-range-int", 70000, true},
		{"int-slice", []interface{}{80, 443}, false},
		{"range-string", "1000-2000", false},
		{"invalid-range-string", "2000-1000", true},
		{"invalid-port-string", "abc", true},
		{"mixed-slice", []interface{}{80, "1000-2000"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePorts(tt.ports)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePorts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
