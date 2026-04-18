package backend

import (
	"context"
	"fmt"
	"sync"

	"github.com/rampartfw/rampart/internal/model"
)

// MultiBackend fans out all operations to multiple underlying backends.
type MultiBackend struct {
	backends []Backend
}

func NewMultiBackend(backends ...Backend) *MultiBackend {
	return &MultiBackend{backends: backends}
}

func (m *MultiBackend) Name() string {
	return "multi-orchestrator"
}

func (m *MultiBackend) Capabilities() model.BackendCapabilities {
	// Intersection of capabilities: we only claim what ALL backends can do
	if len(m.backends) == 0 {
		return model.BackendCapabilities{}
	}
	
	res := m.backends[0].Capabilities()
	for _, b := range m.backends[1:] {
		caps := b.Capabilities()
		res.IPv4 = res.IPv4 && caps.IPv4
		res.IPv6 = res.IPv6 && caps.IPv6
		res.AtomicReplace = res.AtomicReplace && caps.AtomicReplace
		// ... and so on
	}
	return res
}

func (m *MultiBackend) Probe() error {
	for _, b := range m.backends {
		if err := b.Probe(); err != nil {
			return fmt.Errorf("backend %s probe failed: %w", b.Name(), err)
		}
	}
	return nil
}

func (m *MultiBackend) CurrentState(ctx context.Context) (*model.CompiledRuleSet, error) {
	// Usually we return the state of the first backend as "source of truth"
	if len(m.backends) == 0 {
		return &model.CompiledRuleSet{}, nil
	}
	return m.backends[0].CurrentState(ctx)
}

func (m *MultiBackend) Apply(ctx context.Context, rs *model.CompiledRuleSet) error {
	var wg sync.WaitGroup
	errs := make(chan error, len(m.backends))

	for _, b := range m.backends {
		wg.Add(1)
		go func(target Backend) {
			defer wg.Done()
			if err := target.Apply(ctx, rs); err != nil {
				errs <- fmt.Errorf("%s: %v", target.Name(), err)
			}
		}(b)
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		var msg string
		for err := range errs {
			msg += err.Error() + "; "
		}
		return fmt.Errorf("multi-apply failed: %s", msg)
	}
	return nil
}

func (m *MultiBackend) DryRun(ctx context.Context, rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	// Return the most complex plan or the first one
	return m.backends[0].DryRun(ctx, rs)
}

func (m *MultiBackend) Rollback(ctx context.Context, snapshot *model.Snapshot) error {
	for _, b := range m.backends {
		if err := b.Rollback(ctx, snapshot); err != nil {
			return err
		}
	}
	return nil
}

func (m *MultiBackend) Flush(ctx context.Context) error {
	for _, b := range m.backends {
		if err := b.Flush(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *MultiBackend) Stats(ctx context.Context) (map[string]model.RuleStats, error) {
	// Aggregate stats from all backends
	allStats := make(map[string]model.RuleStats)
	for _, b := range m.backends {
		stats, err := b.Stats(ctx)
		if err != nil {
			continue
		}
		for id, s := range stats {
			combined := allStats[id]
			combined.Packets += s.Packets
			combined.Bytes += s.Bytes
			allStats[id] = combined
		}
	}
	return allStats, nil
}

func (m *MultiBackend) Close() error {
	for _, b := range m.backends {
		b.Close()
	}
	return nil
}
