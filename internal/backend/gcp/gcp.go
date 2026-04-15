package gcp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

type GCPBackend struct {
	cfg       backend.BackendConfig
	projectID string
	network   string
	keyFile   string
	client    *http.Client
}

func init() {
	backend.Register("gcp", func(cfg backend.BackendConfig) (backend.Backend, error) {
		return &GCPBackend{
			cfg:       cfg,
			projectID: cfg.Settings["projectId"],
			network:   cfg.Settings["network"],
			keyFile:   cfg.Settings["keyFile"],
			client:    &http.Client{Timeout: 30 * time.Second},
		}, nil
	})
}

func (b *GCPBackend) Name() string {
	return "gcp"
}

func (b *GCPBackend) Capabilities() model.BackendCapabilities {
	return model.BackendCapabilities{
		IPv4:               true,
		IPv6:               true,
		ConnectionTracking: true,
		AtomicReplace:      false,
	}
}

func (b *GCPBackend) Probe() error {
	if b.projectID == "" || b.network == "" || b.keyFile == "" {
		return fmt.Errorf("missing GCP credentials or configuration")
	}
	return nil
}

func (b *GCPBackend) CurrentState(ctx context.Context) (*model.CompiledRuleSet, error) {
	// List GCP Firewall Rules
	return &model.CompiledRuleSet{
		Rules: []model.CompiledRule{},
	}, nil
}

func (b *GCPBackend) Apply(ctx context.Context, rs *model.CompiledRuleSet) error {
	// 1. Get current rules
	// 2. Diff
	// 3. Create/Delete rules
	return nil
}

func (b *GCPBackend) DryRun(ctx context.Context, rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return &model.ExecutionPlan{
		PlannedRuleCount: len(rs.Rules),
	}, nil
}

func (b *GCPBackend) Rollback(ctx context.Context, snapshot *model.Snapshot) error {
	return fmt.Errorf("rollback not implemented for gcp")
}

func (b *GCPBackend) Flush(ctx context.Context) error {
	return nil
}

func (b *GCPBackend) Stats(ctx context.Context) (map[string]model.RuleStats, error) {
	return nil, nil // GCP doesn't provide per-rule stats in this API
}

func (b *GCPBackend) Close() error {
	return nil
}
