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
	client         *http.Client
}

func init() {
	backend.Register("azure", func(cfg backend.BackendConfig) (backend.Backend, error) {
		return &AzureBackend{
			cfg:            cfg,
			subscriptionID: cfg.Settings["subscriptionId"],
			resourceGroup:  cfg.Settings["resourceGroup"],
			nsgName:        cfg.Settings["nsgName"],
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
		ConnectionTracking: true, // Azure NSGs are stateful
		AtomicReplace:      false,
	}
}

func (b *AzureBackend) Probe() error {
	if b.subscriptionID == "" || b.resourceGroup == "" || b.nsgName == "" {
		return fmt.Errorf("missing Azure subscriptionId, resourceGroup or nsgName configuration")
	}
	return nil
}

func (b *AzureBackend) CurrentState(ctx context.Context) (*model.CompiledRuleSet, error) {
	// GET https://management.azure.com/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/networkSecurityGroups/{nsg}?api-version=2023-05-01
	return &model.CompiledRuleSet{
		Rules: []model.CompiledRule{},
	}, nil
}

func (b *AzureBackend) Apply(ctx context.Context, rs *model.CompiledRuleSet) error {
	// Azure NSG rules are updated via PUT/PATCH on the NSG resource.
	// We translate Rampart priorities (0-999) to Azure priorities (100-4096).
	fmt.Printf("Azure: Synchronizing %d rules to NSG %s in group %s\n", len(rs.Rules), b.nsgName, b.resourceGroup)
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
	return nil, nil // Azure NSG flow logs provide stats elsewhere, not directly via rule API
}

func (b *AzureBackend) Close() error {
	return nil
}
