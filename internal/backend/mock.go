package backend

import (
	"github.com/rampartfw/rampart/internal/model"
)

type MockBackend struct {
	Rules   []model.CompiledRule
	Applied int
	ProbeOK bool
}

func (m *MockBackend) Name() string { return "mock" }

func (m *MockBackend) Capabilities() model.BackendCapabilities {
	return model.BackendCapabilities{
		IPv4:          true,
		IPv6:          true,
		AtomicReplace: true,
	}
}

func (m *MockBackend) Probe() error {
	if !m.ProbeOK {
		return nil
	}
	return nil
}

func (m *MockBackend) CurrentState() (*model.CompiledRuleSet, error) {
	return &model.CompiledRuleSet{Rules: m.Rules}, nil
}

func (m *MockBackend) Apply(rs *model.CompiledRuleSet) error {
	m.Rules = rs.Rules
	m.Applied++
	return nil
}

func (m *MockBackend) DryRun(rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return &model.ExecutionPlan{PlannedRuleCount: len(rs.Rules)}, nil
}

func (m *MockBackend) Rollback(snapshot *model.Snapshot) error {
	return nil
}

func (m *MockBackend) Flush() error {
	m.Rules = nil
	return nil
}

func (m *MockBackend) Stats() (map[string]model.RuleStats, error) {
	return make(map[string]model.RuleStats), nil
}

func (m *MockBackend) Close() error {
	return nil
}
