package gcp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

type GCPBackend struct {
	cfg        backend.BackendConfig
	projectID  string
	network    string
	keyFile    string
	accessToken string
	expiry      time.Time
}

func init() {
	backend.Register("gcp", func(cfg backend.BackendConfig) (backend.Backend, error) {
		return &GCPBackend{
			cfg:       cfg,
			projectID: cfg.Settings["projectId"],
			network:   cfg.Settings["network"],
			keyFile:   cfg.Settings["keyFile"],
		}, nil
	})
}

func (b *GCPBackend) Name() string {
	return "gcp"
}

func (b *GCPBackend) Capabilities() model.BackendCapabilities {
	return model.BackendCapabilities{
		IPv4:               true,
		IPv6:               true,
		ConnectionTracking: true,
		AtomicReplace:      false,
	}
}

func (b *GCPBackend) Probe() error {
	if b.projectID == "" || b.network == "" || b.keyFile == "" {
		return fmt.Errorf("missing GCP credentials or configuration")
	}
	return nil
}

func (b *GCPBackend) CurrentState() (*model.CompiledRuleSet, error) {
	// List GCP Firewall Rules
	return &model.CompiledRuleSet{
		Rules: []model.CompiledRule{},
	}, nil
}

func (b *GCPBackend) Apply(rs *model.CompiledRuleSet) error {
	// 1. Get current rules
	// 2. Diff
	// 3. Create/Delete rules
	return nil
}

func (b *GCPBackend) getAccessToken() (string, error) {
	if b.accessToken != "" && time.Now().Before(b.expiry) {
		return b.accessToken, nil
	}

	// 1. Load service account JSON key
	// In a real implementation, we'd parse b.keyFile.
	// 2. Create JWT assertion
	// 3. POST https://oauth2.googleapis.com/token
	
	// Simplified OAuth2 logic placeholder
	return "placeholder-token", nil
}

func (b *GCPBackend) call(method, url string, body []byte) ([]byte, error) {
	token, err := b.getAccessToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GCP API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (b *GCPBackend) DryRun(rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return &model.ExecutionPlan{
		PlannedRuleCount: len(rs.Rules),
	}, nil
}

func (b *GCPBackend) Rollback(snapshot *model.Snapshot) error {
	return fmt.Errorf("rollback not implemented for gcp")
}

func (b *GCPBackend) Flush() error {
	return nil
}

func (b *GCPBackend) Stats() (map[string]model.RuleStats, error) {
	return nil, nil // GCP doesn't provide per-rule stats in this API
}

func (b *GCPBackend) Close() error {
	return nil
}
