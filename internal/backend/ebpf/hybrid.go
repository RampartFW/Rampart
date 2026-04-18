package ebpf

import (
	"context"
	"fmt"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

// HybridBackend combines eBPF for fast path and another backend (e.g., nftables) for complex rules.
type HybridBackend struct {
	fastPath backend.Backend
	slowPath backend.Backend
}

func NewHybridBackend(fast, slow backend.Backend) *HybridBackend {
	return &HybridBackend{
		fastPath: fast,
		slowPath: slow,
	}
}

func init() {
	backend.Register("hybrid", func(cfg backend.BackendConfig) (backend.Backend, error) {
		nft, err := backend.NewBackend("nftables", backend.BackendConfig{Type: "nftables"})
		if err != nil {
			return nil, err
		}
		ebpf, err := backend.NewBackend("ebpf", backend.BackendConfig{Type: "ebpf"})
		if err != nil {
			return nil, err
		}
		return NewHybridBackend(ebpf, nft), nil
	})
}

func (b *HybridBackend) Name() string {
	return fmt.Sprintf("hybrid(%s+%s)", b.fastPath.Name(), b.slowPath.Name())
}

func (b *HybridBackend) Capabilities() model.BackendCapabilities {
	// Hybrid backend combines capabilities, preferring the most capable
	caps := b.slowPath.Capabilities()
	fastCaps := b.fastPath.Capabilities()
	
	if fastCaps.PerRuleCounters {
		caps.PerRuleCounters = true
	}
	return caps
}

func (b *HybridBackend) Probe() error {
	if err := b.fastPath.Probe(); err != nil {
		return fmt.Errorf("fast path probe failed: %w", err)
	}
	return b.slowPath.Probe()
}

func (b *HybridBackend) CurrentState(ctx context.Context) (*model.CompiledRuleSet, error) {
	// Combined state from both backends
	slowState, err := b.slowPath.CurrentState(ctx)
	if err != nil {
		return nil, err
	}
	fastState, err := b.fastPath.CurrentState(ctx)
	if err != nil {
		return nil, err
	}

	slowState.Rules = append(slowState.Rules, fastState.Rules...)
	return slowState, nil
}

func (b *HybridBackend) Apply(ctx context.Context, rs *model.CompiledRuleSet) error {
	// Split rules into fast path and slow path
	var fastRules []model.CompiledRule
	var slowRules []model.CompiledRule

	for _, rule := range rs.Rules {
		if b.isFastPathCapable(rule) {
			fastRules = append(fastRules, rule)
		} else {
			slowRules = append(slowRules, rule)
		}
	}

	// Apply to both backends
	if err := b.fastPath.Apply(ctx, &model.CompiledRuleSet{Rules: fastRules, Metadata: rs.Metadata}); err != nil {
		return fmt.Errorf("fast path apply failed: %w", err)
	}
	return b.slowPath.Apply(ctx, &model.CompiledRuleSet{Rules: slowRules, Metadata: rs.Metadata})
}

func (b *HybridBackend) isFastPathCapable(r model.CompiledRule) bool {
	// eBPF/XDP backend currently only supports simple CIDR/Port matching
	// No connection tracking or complex schedules
	if r.Schedule != nil {
		return false
	}
	if len(r.Match.States) > 0 {
		return false
	}
	return true
}

func (b *HybridBackend) DryRun(ctx context.Context, rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return b.slowPath.DryRun(ctx, rs)
}

func (b *HybridBackend) Rollback(ctx context.Context, snapshot *model.Snapshot) error {
	// Simple rollback strategy: apply through the hybrid interface
	var rs model.CompiledRuleSet
	// This would require unmarshaling the snapshot rules
	return b.Apply(ctx, &rs)
}

func (b *HybridBackend) Flush(ctx context.Context) error {
	if err := b.fastPath.Flush(ctx); err != nil {
		return err
	}
	return b.slowPath.Flush(ctx)
}

func (b *HybridBackend) Stats(ctx context.Context) (map[string]model.RuleStats, error) {
	fastStats, err := b.fastPath.Stats(ctx)
	if err != nil {
		return nil, err
	}
	slowStats, err := b.slowPath.Stats(ctx)
	if err != nil {
		return nil, err
	}

	// Merge stats
	for k, v := range fastStats {
		slowStats[k] = v
	}
	return slowStats, nil
}

func (b *HybridBackend) Close() error {
	b.fastPath.Close()
	return b.slowPath.Close()
}
