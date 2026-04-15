package azure

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

type AzureBackend struct {
	cfg            backend.BackendConfig
	subscriptionID string
	resourceGroup  string
	nsgName        string
	tenantID       string
	clientID       string
	clientSecret   string
	accessToken    string
	expiry         time.Time
	client         *http.Client
}

func init() {
	backend.Register("azure", func(cfg backend.BackendConfig) (backend.Backend, error) {
		return &AzureBackend{
			cfg:            cfg,
			subscriptionID: cfg.Settings["subscriptionId"],
			resourceGroup:  cfg.Settings["resourceGroup"],
			nsgName:        cfg.Settings["nsgName"],
			tenantID:       cfg.Settings["tenantId"],
			clientID:       cfg.Settings["clientId"],
			clientSecret:   cfg.Settings["clientSecret"],
			client:         &http.Client{Timeout: 30 * time.Second},
		}, nil
	})
}

func (b *AzureBackend) Name() string {
	return "azure"
}

func (b *AzureBackend) Capabilities() model.BackendCapabilities {
	return model.BackendCapabilities{
		IPv4:               true,
		IPv6:               true,
		ConnectionTracking: true,
		AtomicReplace:      false,
	}
}

func (b *AzureBackend) Probe() error {
	if b.subscriptionID == "" || b.resourceGroup == "" || b.nsgName == "" || b.tenantID == "" || b.clientID == "" || b.clientSecret == "" {
		return fmt.Errorf("missing Azure credentials or configuration")
	}
	return nil
}

func (b *AzureBackend) CurrentState(ctx context.Context) (*model.CompiledRuleSet, error) {
	// List Azure NSG Security Rules
	return &model.CompiledRuleSet{
		Rules: []model.CompiledRule{},
	}, nil
}

func (b *AzureBackend) Apply(ctx context.Context, rs *model.CompiledRuleSet) error {
	// 1. Get current rules
	// 2. Diff
	// 3. Create/Delete rules
	return nil
}

func (b *AzureBackend) DryRun(ctx context.Context, rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return &model.ExecutionPlan{
		PlannedRuleCount: len(rs.Rules),
	}, nil
}

func (b *AzureBackend) Rollback(ctx context.Context, snapshot *model.Snapshot) error {
	return fmt.Errorf("rollback not implemented for azure")
}

func (b *AzureBackend) Flush(ctx context.Context) error {
	return nil
}

func (b *AzureBackend) Stats(ctx context.Context) (map[string]model.RuleStats, error) {
	return nil, nil // Azure NSG doesn't provide real-time per-rule stats in this API
}

func (b *AzureBackend) Close() error {
	return nil
}
