package engine

import (
	"strings"
	"testing"

	"github.com/rampartfw/rampart/internal/model"
)

func TestFormatConflicts(t *testing.T) {
	conflicts := []model.Conflict{
		{
			Type:     model.ConflictShadow,
			Severity: model.SeverityWarning,
			RuleA:    model.CompiledRule{Name: "rule1", Priority: 10, PolicyName: "policy1"},
			RuleB:    model.CompiledRule{Name: "rule2", Priority: 20, PolicyName: "policy1"},
			Message:  "Rule \"rule2\" is shadowed by higher priority rule \"rule1\".",
		},
	}

	t.Run("Text format", func(t *testing.T) {
		formatted := FormatConflicts(conflicts, "text")
		if !strings.Contains(formatted, "Conflict Report") {
			t.Errorf("Expected report header, got: %s", formatted)
		}
		if !strings.Contains(formatted, "rule1") {
			t.Errorf("Expected rule1 name, got: %s", formatted)
		}
		if !strings.Contains(formatted, "shadow") {
			t.Errorf("Expected conflict type, got: %s", formatted)
		}
	})

	t.Run("JSON format", func(t *testing.T) {
		formatted := FormatConflicts(conflicts, "json")
		if !strings.HasPrefix(formatted, "[") {
			t.Errorf("Expected JSON array, got: %s", formatted)
		}
		if !strings.Contains(formatted, "shadow") {
			t.Errorf("Expected shadow type in JSON, got: %s", formatted)
		}
	})
}
