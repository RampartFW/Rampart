package engine

import (
	"os"
	"testing"
)

func TestSubstituteVars(t *testing.T) {
	vars := map[string]interface{}{
		"ssh_port": 22,
		"subnet":   "10.0.1.0/24",
	}

	data := []byte(`
rules:
  - name: allow-ssh
    match:
      protocol: tcp
      destPorts: ["${ssh_port}"]
      sourceCIDRs: ["${subnet}"]
    action: accept
`)

	expected := `
rules:
  - name: allow-ssh
    match:
      protocol: tcp
      destPorts: ["22"]
      sourceCIDRs: ["10.0.1.0/24"]
    action: accept
`

	result, err := SubstituteVars(data, vars)
	if err != nil {
		t.Fatalf("SubstituteVars failed: %v", err)
	}

	if string(result) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(result))
	}
}

func TestSubstituteVars_Unresolved(t *testing.T) {
	vars := map[string]interface{}{}
	data := []byte(`port: ${ssh_port}`)

	_, err := SubstituteVars(data, vars)
	if err == nil {
		t.Error("expected error for unresolved variable, got nil")
	}
}

func TestParseVariablesFile(t *testing.T) {
	content := `
apiVersion: rampart.dev/v1
kind: Variables
metadata:
  name: test-vars
variables:
  ssh_port: 22
  subnet: 10.0.1.0/24
`
	tmpfile, err := os.CreateTemp("", "vars-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	vars, err := ParseVariablesFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseVariablesFile failed: %v", err)
	}

	if vars["ssh_port"] != 22 {
		t.Errorf("expected ssh_port 22, got %v", vars["ssh_port"])
	}
	if vars["subnet"] != "10.0.1.0/24" {
		t.Errorf("expected subnet 10.0.1.0/24, got %v", vars["subnet"])
	}
}
