package aws

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

type AWSBackend struct {
	cfg       backend.BackendConfig
	accessKey string
	secretKey string
	region    string
	sgID      string
	client    *http.Client
}

func init() {
	backend.Register("aws", func(cfg backend.BackendConfig) (backend.Backend, error) {
		return &AWSBackend{
			cfg:       cfg,
			accessKey: cfg.Settings["accessKey"],
			secretKey: cfg.Settings["secretKey"],
			region:    cfg.Settings["region"],
			sgID:      cfg.Settings["securityGroupId"],
			client:    &http.Client{Timeout: 30 * time.Second},
		}, nil
	})
}

func (b *AWSBackend) Name() string {
	return "aws"
}

func (b *AWSBackend) Capabilities() model.BackendCapabilities {
	return model.BackendCapabilities{
		IPv4:               true,
		IPv6:               true,
		ConnectionTracking: true, // AWS SGs are always stateful
		AtomicReplace:      false,
	}
}

func (b *AWSBackend) Probe() error {
	if b.accessKey == "" || b.secretKey == "" || b.region == "" || b.sgID == "" {
		return fmt.Errorf("missing AWS credentials or configuration")
	}
	return nil
}

func (b *AWSBackend) CurrentState(ctx context.Context) (*model.CompiledRuleSet, error) {
	// DescribeSecurityGroupRules
	return &model.CompiledRuleSet{
		Rules: []model.CompiledRule{},
	}, nil
}

func (b *AWSBackend) Apply(ctx context.Context, rs *model.CompiledRuleSet) error {
	// 1. Get current state from AWS API
	current, err := b.CurrentState(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current AWS state: %w", err)
	}

	// 2. Diff: find rules to add and remove
	plan := b.generatePlan(current, rs)
	
	// 3. Apply changes (Authorize and Revoke)
	for _, rule := range plan.Removed {
		if err := b.revokeRule(ctx, rule); err != nil {
			return fmt.Errorf("failed to revoke rule %s: %w", rule.Name, err)
		}
	}
	for _, rule := range plan.Added {
		if err := b.authorizeRule(ctx, rule); err != nil {
			return fmt.Errorf("failed to authorize rule %s: %w", rule.Name, err)
		}
	}

	return nil
}

func (b *AWSBackend) generatePlan(current, desired *model.CompiledRuleSet) *awsPlan {
	plan := &awsPlan{}
	
	// Set for quick lookup
	currentMap := make(map[string]model.CompiledRule)
	for _, r := range current.Rules {
		currentMap[r.Name] = r
	}

	desiredMap := make(map[string]model.CompiledRule)
	for _, r := range desired.Rules {
		desiredMap[r.Name] = r
		if _, ok := currentMap[r.Name]; !ok {
			plan.Added = append(plan.Added, r)
		}
	}

	for _, r := range current.Rules {
		if _, ok := desiredMap[r.Name]; !ok {
			plan.Removed = append(plan.Removed, r)
		}
	}

	return plan
}

type awsPlan struct {
	Added   []model.CompiledRule
	Removed []model.CompiledRule
}

func (b *AWSBackend) authorizeRule(ctx context.Context, rule model.CompiledRule) error {
	// This would construct a SigV4 signed request to EC2 AuthorizeSecurityGroupIngress
	fmt.Printf("AWS: Authorizing rule %s in SG %s\n", rule.Name, b.sgID)
	return nil
}

func (b *AWSBackend) revokeRule(ctx context.Context, rule model.CompiledRule) error {
	// This would construct a SigV4 signed request to EC2 RevokeSecurityGroupIngress
	fmt.Printf("AWS: Revoking rule %s in SG %s\n", rule.Name, b.sgID)
	return nil
}

func (b *AWSBackend) DryRun(ctx context.Context, rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return &model.ExecutionPlan{
		PlannedRuleCount: len(rs.Rules),
	}, nil
}

func (b *AWSBackend) Rollback(ctx context.Context, snapshot *model.Snapshot) error {
	return fmt.Errorf("rollback not implemented for aws")
}

func (b *AWSBackend) Flush(ctx context.Context) error {
	return nil
}

func (b *AWSBackend) Stats(ctx context.Context) (map[string]model.RuleStats, error) {
	return nil, nil // AWS SGs don't provide per-rule stats
}

func (b *AWSBackend) Close() error {
	return nil
}
