//go:build !linux

package ebpf

import (
	"context"
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

func (b *EBPFBackend) CurrentState(ctx context.Context) (*model.CompiledRuleSet, error) {
	return nil, fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) Apply(ctx context.Context, rs *model.CompiledRuleSet) error {
	return fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) DryRun(ctx context.Context, rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return nil, fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) Rollback(ctx context.Context, snapshot *model.Snapshot) error {
	return fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) Flush(ctx context.Context) error {
	return fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) Stats(ctx context.Context) (map[string]model.RuleStats, error) {
	return nil, fmt.Errorf("ebpf backend only supported on Linux")
}

func (b *EBPFBackend) Close() error {
	return nil
}
