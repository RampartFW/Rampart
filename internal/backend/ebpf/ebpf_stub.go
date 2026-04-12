//go:build !linux

package ebpf

import (
	"fmt"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

type EBPFBackend struct {
	cfg backend.BackendConfig
}

func init() {
	backend.Register("ebpf", func(cfg backend.BackendConfig) (backend.Backend, error) {
		return &EBPFBackend{cfg: cfg}, nil
	})
}

func (b *EBPFBackend) Name() string {
	return "ebpf"
}

func (b *EBPFBackend) Capabilities() model.BackendCapabilities {
	return model.BackendCapabilities{}
}

func (b *EBPFBackend) Probe() error {
	return fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) CurrentState() (*model.CompiledRuleSet, error) {
	return nil, fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) Apply(rs *model.CompiledRuleSet) error {
	return fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) DryRun(rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return nil, fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) Rollback(snapshot *model.Snapshot) error {
	return fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) Flush() error {
	return fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) Stats() (map[string]model.RuleStats, error) {
	return nil, fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) Close() error {
	return nil
}
