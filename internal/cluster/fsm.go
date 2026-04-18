package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/model"
)

// PolicyFSM implements the Raft FSM interface for Rampart.
type PolicyFSM struct {
	backend backend.Backend
	engine  *engine.Engine
}

// NewPolicyFSM creates a new PolicyFSM.
func NewPolicyFSM(b backend.Backend, e *engine.Engine) *PolicyFSM {
	return &PolicyFSM{
		backend: b,
		engine:  e,
	}
}

// Apply applies a committed log entry to the local backend.
func (f *PolicyFSM) Apply(entry model.LogEntry) error {
	switch entry.Type {
	case model.EntryPolicyUpdate:
		return f.applyPolicyUpdate(entry.Data)
	case model.EntryIPBan:
		return f.applyIPBan(entry.Data)
	case model.EntryConfigChange:
		// Handle config change
		return nil
	case model.EntryNodeJoin, model.EntryNodeLeave:
		// These are usually handled by the Raft implementation itself
		return nil
	default:
		return fmt.Errorf("unknown entry type: %s", entry.Type)
	}
}

func (f *PolicyFSM) applyIPBan(data []byte) error {
	var ip string
	if err := json.Unmarshal(data, &ip); err != nil {
		return fmt.Errorf("failed to unmarshal IP ban: %w", err)
	}

	// 1. Add to in-memory ruleset as a high-priority block rule
	current := f.engine.CurrentRules()
	if current == nil {
		current = &model.CompiledRuleSet{}
	}

	// Create a dynamic block rule
	banRule := model.CompiledRule{
		ID:        fmt.Sprintf("ips-ban-%s", ip),
		Name:      fmt.Sprintf("IPS Auto-Ban: %s", ip),
		Priority:  -100, // Very high priority
		Action:    model.ActionDrop,
		Direction: model.DirectionInbound,
		Match: model.CompiledMatch{
			SourceNets: []net.IPNet{
				{IP: net.ParseIP(ip), Mask: net.CIDRMask(32, 32)},
			},
		},
	}

	current.Rules = append([]model.CompiledRule{banRule}, current.Rules...)
	f.engine.SetRules(current)

	// 2. Re-apply to backend
	return f.engine.ReapplyRules(context.Background())
}

func (f *PolicyFSM) applyPolicyUpdate(data []byte) error {
	var ps model.PolicySetYAML
	if err := json.Unmarshal(data, &ps); err != nil {
		return fmt.Errorf("failed to unmarshal policy: %w", err)
	}

	// Compile the policy
	compiled, err := engine.Compile(&ps, nil)
	if err != nil {
		return fmt.Errorf("failed to compile policy: %w", err)
	}

	// Update the engine's in-memory ruleset
	f.engine.SetRules(compiled)

	// Apply to backend
	if err := f.engine.ReapplyRules(context.Background()); err != nil {
		return fmt.Errorf("failed to apply policy to backend: %w", err)
	}

	return nil
}

// Snapshot returns a snapshot of the current state.
func (f *PolicyFSM) Snapshot() ([]byte, error) {
	// For Rampart, the state is the current compiled ruleset
	state, err := f.backend.CurrentState(context.Background())
	if err != nil {
		return nil, err
	}
	return json.Marshal(state)
}

// Restore restores the state from a snapshot.
func (f *PolicyFSM) Restore(snapshot []byte) error {
	var rs model.CompiledRuleSet
	if err := json.Unmarshal(snapshot, &rs); err != nil {
		return err
	}
	return f.backend.Apply(context.Background(), &rs)
}
