package backend

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rampartfw/rampart/internal/model"
)

type MockBackend struct {
	Rules    []model.CompiledRule
	Applied  int
	ProbeOK  bool
	stateDir string
}

func init() {
	Register("mock", func(cfg BackendConfig) (Backend, error) {
		m := &MockBackend{ProbeOK: true, stateDir: "./tmp/mock_state"}
		_ = os.MkdirAll(m.stateDir, 0755)
		_ = m.load()
		return m, nil
	})
}

func (m *MockBackend) load() error {
	data, err := os.ReadFile(filepath.Join(m.stateDir, "rules.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &m.Rules)
}

func (m *MockBackend) save() error {
	data, _ := json.Marshal(m.Rules)
	return os.WriteFile(filepath.Join(m.stateDir, "rules.json"), data, 0644)
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
	return nil
}

func (m *MockBackend) CurrentState(ctx context.Context) (*model.CompiledRuleSet, error) {
	_ = m.load() // Always reload for CLI
	return &model.CompiledRuleSet{Rules: m.Rules}, nil
}

func (m *MockBackend) Apply(ctx context.Context, rs *model.CompiledRuleSet) error {
	m.Rules = rs.Rules
	m.Applied++
	return m.save()
}

func (m *MockBackend) DryRun(ctx context.Context, rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return &model.ExecutionPlan{PlannedRuleCount: len(rs.Rules)}, nil
}

func (m *MockBackend) Rollback(ctx context.Context, snapshot *model.Snapshot) error {
	return nil
}

func (m *MockBackend) Flush(ctx context.Context) error {
	m.Rules = nil
	return m.save()
}

func (m *MockBackend) Stats(ctx context.Context) (map[string]model.RuleStats, error) {
	return make(map[string]model.RuleStats), nil
}

func (m *MockBackend) Close() error {
	return nil
}
