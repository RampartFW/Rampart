package engine

import (
	"os"
	"testing"
)

func TestParsePolicyFile(t *testing.T) {
	content := `
apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: test-policy
policies:
  - name: p1
    priority: 100
    rules:
      - name: r1
        match:
          protocol: tcp
          destPorts: 80
        action: accept
`
	tmpfile, err := os.CreateTemp("", "policy-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	ps, err := ParsePolicyFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParsePolicyFile failed: %v", err)
	}

	if ps.Metadata.Name != "test-policy" {
		t.Errorf("expected name test-policy, got %s", ps.Metadata.Name)
	}

	if len(ps.Policies) != 1 {
		t.Errorf("expected 1 policy, got %d", len(ps.Policies))
	}

	if ps.Policies[0].Name != "p1" {
		t.Errorf("expected policy name p1, got %s", ps.Policies[0].Name)
	}
}

func TestParsePolicyFileWithVars(t *testing.T) {
	content := `
apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: test-policy
policies:
  - name: p1
    rules:
      - name: r1
        match:
          protocol: tcp
          destPorts: ${web_port}
        action: accept
`
	tmpfile, err := os.CreateTemp("", "policy-vars-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	vars := map[string]interface{}{
		"web_port": 80,
	}

	ps, err := ParsePolicyFileWithVars(tmpfile.Name(), vars)
	if err != nil {
		t.Fatalf("ParsePolicyFileWithVars failed: %v", err)
	}

	if ps.Policies[0].Rules[0].Match.DestPorts != 80 {
		t.Errorf("expected destPorts 80, got %v", ps.Policies[0].Rules[0].Match.DestPorts)
	}
}

func TestParsePolicyData_CIDRNormalization(t *testing.T) {
	content := `
apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: test-policy
policies:
  - name: p1
    rules:
      - name: r1
        match:
          sourceCIDRs: ["10.0.1.1", "2001:db8::1", "192.168.1.0/24"]
        action: accept
`
	ps, err := ParsePolicyData([]byte(content))
	if err != nil {
		t.Fatalf("ParsePolicyData failed: %v", err)
	}

	cidrs := ps.Policies[0].Rules[0].Match.SourceCIDRs
	if len(cidrs) != 3 {
		t.Fatalf("expected 3 CIDRs, got %d", len(cidrs))
	}

	if cidrs[0] != "10.0.1.1/32" {
		t.Errorf("expected 10.0.1.1/32, got %s", cidrs[0])
	}
	if cidrs[1] != "2001:db8::1/128" {
		t.Errorf("expected 2001:db8::1/128, got %s", cidrs[1])
	}
	if cidrs[2] != "192.168.1.0/24" {
		t.Errorf("expected 192.168.1.0/24, got %s", cidrs[2])
	}
}

