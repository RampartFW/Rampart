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
	// 1. Get current rules
	// 2. Diff
	// 3. Authorize/Revoke
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
