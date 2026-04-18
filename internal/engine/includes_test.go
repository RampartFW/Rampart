package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rampartfw/rampart/internal/model"
)

func TestResolveIncludes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "includes-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	basePath := filepath.Join(tmpDir, "base.yaml")
	incPath := filepath.Join(tmpDir, "inc.yaml")

	baseContent := `
apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: base
includes:
  - path: inc.yaml
policies:
  - name: base-policy
    rules:
      - name: base-rule
        action: accept
`
	incContent := `
apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: inc
policies:
  - name: inc-policy
    rules:
      - name: inc-rule
        action: drop
`

	if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(incPath, []byte(incContent), 0644); err != nil {
		t.Fatal(err)
	}

	ps, err := ParsePolicyFile(basePath)
	if err != nil {
		t.Fatalf("ParsePolicyFile failed: %v", err)
	}

	if err := ResolveIncludes(ps, basePath); err != nil {
		t.Fatalf("ResolveIncludes failed: %v", err)
	}

	if len(ps.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(ps.Policies))
	}

	policyNames := make(map[string]bool)
	for _, p := range ps.Policies {
		policyNames[p.Name] = true
	}

	if !policyNames["base-policy"] {
		t.Errorf("expected policy base-policy missing")
	}
	if !policyNames["inc-policy"] {
		t.Errorf("expected policy inc-policy missing")
	}
}

func TestResolveIncludes_Circular(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "includes-circular")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	p1 := filepath.Join(tmpDir, "p1.yaml")
	p2 := filepath.Join(tmpDir, "p2.yaml")

	p1Content := `
apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: p1
includes:
  - path: p2.yaml
`
	p2Content := `
apiVersion: rampartfw.com/v1
kind: PolicySet
metadata:
  name: p2
includes:
  - path: p1.yaml
`

	if err := os.WriteFile(p1, []byte(p1Content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p2, []byte(p2Content), 0644); err != nil {
		t.Fatal(err)
	}

	ps, _ := ParsePolicyFile(p1)
	err = ResolveIncludes(ps, p1)
	if err == nil {
		t.Error("expected circular include error, got nil")
	}
}

func TestResolveIncludes_NotFound(t *testing.T) {
	ps := &model.PolicySetYAML{
		Includes: []model.IncludeRef{
			{Path: "nonexistent.yaml"},
		},
	}

	err := ResolveIncludes(ps, "anywhere/base.yaml")
	if err == nil {
		t.Error("expected error for missing include, got nil")
	}
}

