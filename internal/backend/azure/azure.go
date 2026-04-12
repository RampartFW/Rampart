package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

type AzureBackend struct {
	cfg            backend.BackendConfig
	subscriptionID string
	resourceGroup  string
	nsgName        string
	tenantID       string
	clientID       string
	clientSecret   string
	accessToken    string
	expiry         time.Time
}

func init() {
	backend.Register("azure", func(cfg backend.BackendConfig) (backend.Backend, error) {
		return &AzureBackend{
			cfg:            cfg,
			subscriptionID: cfg.Settings["subscriptionId"],
			resourceGroup:  cfg.Settings["resourceGroup"],
			nsgName:        cfg.Settings["nsgName"],
			tenantID:       cfg.Settings["tenantId"],
			clientID:       cfg.Settings["clientId"],
			clientSecret:   cfg.Settings["clientSecret"],
		}, nil
	})
}

func (b *AzureBackend) Name() string {
	return "azure"
}

func (b *AzureBackend) Capabilities() model.BackendCapabilities {
	return model.BackendCapabilities{
		IPv4:               true,
		IPv6:               true,
		ConnectionTracking: true,
		AtomicReplace:      false,
	}
}

func (b *AzureBackend) Probe() error {
	if b.subscriptionID == "" || b.resourceGroup == "" || b.nsgName == "" || b.tenantID == "" || b.clientID == "" || b.clientSecret == "" {
		return fmt.Errorf("missing Azure credentials or configuration")
	}
	return nil
}

func (b *AzureBackend) CurrentState() (*model.CompiledRuleSet, error) {
	// List Azure NSG Security Rules
	return &model.CompiledRuleSet{
		Rules: []model.CompiledRule{},
	}, nil
}

func (b *AzureBackend) Apply(rs *model.CompiledRuleSet) error {
	// 1. Get current rules
	// 2. Diff
	// 3. Create/Delete rules
	return nil
}

func (b *AzureBackend) getAccessToken() (string, error) {
	if b.accessToken != "" && time.Now().Before(b.expiry) {
		return b.accessToken, nil
	}

	authURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", b.tenantID)
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", b.clientID)
	data.Set("client_secret", b.clientSecret)
	data.Set("resource", "https://management.azure.com/")

	req, err := http.NewRequest("POST", authURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Azure OAuth2 error (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	b.accessToken = result.AccessToken
	b.expiry = time.Now().Add(50 * time.Minute) // Placeholder expiry
	return b.accessToken, nil
}

func (b *AzureBackend) call(method, relativeURL string, body []byte) ([]byte, error) {
	token, err := b.getAccessToken()
	if err != nil {
		return nil, err
	}

	baseURL := "https://management.azure.com"
	fullURL := baseURL + relativeURL

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(body))
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Azure API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (b *AzureBackend) DryRun(rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return &model.ExecutionPlan{
		PlannedRuleCount: len(rs.Rules),
	}, nil
}

func (b *AzureBackend) Rollback(snapshot *model.Snapshot) error {
	return fmt.Errorf("rollback not implemented for azure")
}

func (b *AzureBackend) Flush() error {
	return nil
}

func (b *AzureBackend) Stats() (map[string]model.RuleStats, error) {
	return nil, nil // Azure NSG doesn't provide real-time per-rule stats in this API
}

func (b *AzureBackend) Close() error {
	return nil
}
