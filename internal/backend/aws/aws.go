package aws

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

type AWSBackend struct {
	cfg       backend.BackendConfig
	accessKey string
	secretKey string
	region    string
	sgID      string
}

func init() {
	backend.Register("aws", func(cfg backend.BackendConfig) (backend.Backend, error) {
		return &AWSBackend{
			cfg:       cfg,
			accessKey: cfg.Settings["accessKey"],
			secretKey: cfg.Settings["secretKey"],
			region:    cfg.Settings["region"],
			sgID:      cfg.Settings["securityGroupId"],
		}, nil
	})
}

func (b *AWSBackend) Name() string {
	return "aws"
}

func (b *AWSBackend) Capabilities() model.BackendCapabilities {
	return model.BackendCapabilities{
		IPv4:               true,
		IPv6:               true,
		ConnectionTracking: true, // AWS SGs are always stateful
		AtomicReplace:      false,
	}
}

func (b *AWSBackend) Probe() error {
	if b.accessKey == "" || b.secretKey == "" || b.region == "" || b.sgID == "" {
		return fmt.Errorf("missing AWS credentials or configuration")
	}
	return nil
}

func (b *AWSBackend) CurrentState() (*model.CompiledRuleSet, error) {
	// DescribeSecurityGroupRules
	return &model.CompiledRuleSet{
		Rules: []model.CompiledRule{},
	}, nil
}

func (b *AWSBackend) Apply(rs *model.CompiledRuleSet) error {
	// 1. Get current rules
	// 2. Diff
	// 3. Authorize/Revoke
	return nil
}

func (b *AWSBackend) call(action string, params map[string]string) ([]byte, error) {
	endpoint := fmt.Sprintf("https://ec2.%s.amazonaws.com/", b.region)
	
	// Create body
	body := fmt.Sprintf("Action=%s&Version=2016-11-15", action)
	for k, v := range params {
		body += fmt.Sprintf("&%s=%s", k, v)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	Sign(req, []byte(body), b.accessKey, b.secretKey, b.region, "ec2", time.Now())

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
		return nil, fmt.Errorf("AWS API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (b *AWSBackend) DryRun(rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return &model.ExecutionPlan{
		PlannedRuleCount: len(rs.Rules),
	}, nil
}

func (b *AWSBackend) Rollback(snapshot *model.Snapshot) error {
	return fmt.Errorf("rollback not implemented for aws")
}

func (b *AWSBackend) Flush() error {
	return nil
}

func (b *AWSBackend) Stats() (map[string]model.RuleStats, error) {
	return nil, nil // AWS SGs don't provide per-rule stats
}

func (b *AWSBackend) Close() error {
	return nil
}
