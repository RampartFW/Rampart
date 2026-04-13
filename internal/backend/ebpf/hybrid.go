package ebpf

import (
	"context"
	"fmt"
	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

type HybridBackend struct {
	cfg      backend.BackendConfig
	fastPath backend.Backend
	slowPath backend.Backend
}

func init() {
	backend.Register("hybrid", func(cfg backend.BackendConfig) (backend.Backend, error) {
		fastType, ok := cfg.Settings["fastPath"]
		if !ok {
			fastType = "ebpf"
		}
		slowType, ok := cfg.Settings["slowPath"]
		if !ok {
			slowType = "nftables"
		}

		fast, err := backend.NewBackend(fastType, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create fast path backend %s: %w", fastType, err)
		}
		slow, err := backend.NewBackend(slowType, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create slow path backend %s: %w", slowType, err)
		}

		return &HybridBackend{
			cfg:      cfg,
			fastPath: fast,
			slowPath: slow,
		}, nil
	})
}

func (b *HybridBackend) Name() string {
	return "hybrid"
}

func (b *HybridBackend) Capabilities() model.BackendCapabilities {
	// Intersection or union? Usually union for capabilities.
	fastCaps := b.fastPath.Capabilities()
	slowCaps := b.slowPath.Capabilities()

	return model.BackendCapabilities{
		IPv4:               fastCaps.IPv4 || slowCaps.IPv4,
		IPv6:               fastCaps.IPv6 || slowCaps.IPv6,
		RateLimiting:       fastCaps.RateLimiting || slowCaps.RateLimiting,
		ConnectionTracking: fastCaps.ConnectionTracking || slowCaps.ConnectionTracking,
		Logging:            fastCaps.Logging || slowCaps.Logging,
		PerRuleCounters:    fastCaps.PerRuleCounters || slowCaps.PerRuleCounters,
		AtomicReplace:      fastCaps.AtomicReplace && slowCaps.AtomicReplace,
		InterfaceFiltering: fastCaps.InterfaceFiltering || slowCaps.InterfaceFiltering,
	}
}

func (b *HybridBackend) Probe() error {
	if err := b.fastPath.Probe(); err != nil {
		return fmt.Errorf("fast path probe failed: %w", err)
	}
	if err := b.slowPath.Probe(); err != nil {
		return fmt.Errorf("slow path probe failed: %w", err)
	}
	return nil
}

func (b *HybridBackend) CurrentState(ctx context.Context) (*model.CompiledRuleSet, error) {
	// Reconstruct state from both backends.
	// For now, just return slow path as source of truth.
	return b.slowPath.CurrentState(ctx)
}

func (b *HybridBackend) Apply(ctx context.Context, rs *model.CompiledRuleSet) error {
	// Split rules between fast and slow path
	fastRules := &model.CompiledRuleSet{
		Metadata: rs.Metadata,
		Hash:     rs.Hash,
	}
	slowRules := &model.CompiledRuleSet{
		Metadata: rs.Metadata,
		Hash:     rs.Hash,
	}

	for _, rule := range rs.Rules {
		if b.isFastPathEligible(rule) {
			fastRules.Rules = append(fastRules.Rules, rule)
		} else {
			slowRules.Rules = append(slowRules.Rules, rule)
		}
	}

	// Apply to both
	if err := b.fastPath.Apply(ctx, fastRules); err != nil {
		return fmt.Errorf("fast path apply failed: %w", err)
	}
	if err := b.slowPath.Apply(ctx, slowRules); err != nil {
		return fmt.Errorf("slow path apply failed: %w", err)
	}

	return nil
}

func (b *HybridBackend) isFastPathEligible(rule model.CompiledRule) bool {
	// Simple eligibility logic: system rules (priority < 100) or rate-limit rules
	if rule.Priority < 100 {
		return true
	}
	if rule.Action == model.ActionRateLimit {
		return true
	}
	if rule.Tags["fastpath"] == "true" {
		return true
	}
	return false
}

func (b *HybridBackend) DryRun(ctx context.Context, rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return b.slowPath.DryRun(ctx, rs)
}

func (b *HybridBackend) Rollback(ctx context.Context, snapshot *model.Snapshot) error {
	if err := b.fastPath.Rollback(ctx, snapshot); err != nil {
		return err
	}
	return b.slowPath.Rollback(ctx, snapshot)
}

func (b *HybridBackend) Flush(ctx context.Context) error {
	if err := b.fastPath.Flush(ctx); err != nil {
		return err
	}
	return b.slowPath.Flush(ctx)
}

func (b *HybridBackend) Stats(ctx context.Context) (map[string]model.RuleStats, error) {
	fastStats, _ := b.fastPath.Stats(ctx)
	slowStats, _ := b.slowPath.Stats(ctx)

	// Merge stats
	allStats := make(map[string]model.RuleStats)
	for id, s := range fastStats {
		allStats[id] = s
	}
	for id, s := range slowStats {
		if existing, ok := allStats[id]; ok {
			existing.Packets += s.Packets
			existing.Bytes += s.Bytes
			allStats[id] = existing
		} else {
			allStats[id] = s
		}
	}
	return allStats, nil
}

func (b *HybridBackend) Close() error {
	b.fastPath.Close()
	b.slowPath.Close()
	return nil
}
