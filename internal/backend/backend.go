package backend

import (
	"github.com/rampartfw/rampart/internal/model"
)

// Backend is the core interface that all firewall backends must implement.
type Backend interface {
	// Name returns the backend identifier (e.g., "nftables", "iptables", "ebpf")
	Name() string

	// Capabilities reports what this backend supports
	Capabilities() model.BackendCapabilities

	// Probe checks if this backend is available on the current system
	Probe() error

	// CurrentState returns the active firewall rules in normalized form
	CurrentState() (*model.CompiledRuleSet, error)

	// Apply atomically applies a complete RuleSet, replacing all managed rules
	Apply(rs *model.CompiledRuleSet) error

	// DryRun returns what Apply would do without actually doing it
	DryRun(rs *model.CompiledRuleSet) (*model.ExecutionPlan, error)

	// Rollback restores a previously captured snapshot
	Rollback(snapshot *model.Snapshot) error

	// Flush removes all Rampart-managed rules (leaves system rules intact)
	Flush() error

	// Stats returns per-rule packet/byte counters
	Stats() (map[string]model.RuleStats, error)

	// Close releases any resources held by the backend
	Close() error
}

// BackendConfig holds the configuration for a backend
type BackendConfig struct {
	Type     string
	Settings map[string]string
}
