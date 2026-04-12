package engine

import (
	"testing"

	"github.com/rampartfw/rampart/internal/model"
)

func TestGeneratePlan(t *testing.T) {
	current := &model.CompiledRuleSet{
		Rules: []model.CompiledRule{
			{
				ID:         "1",
				Name:       "rule1",
				PolicyName: "policy1",
				Action:     model.ActionAccept,
			},
			{
				ID:         "2",
				Name:       "rule2",
				PolicyName: "policy1",
				Action:     model.ActionDrop,
			},
		},
	}

	desired := &model.CompiledRuleSet{
		Rules: []model.CompiledRule{
			{
				ID:         "1", // No change
				Name:       "rule1",
				PolicyName: "policy1",
				Action:     model.ActionAccept,
			},
			{
				ID:         "3", // Modified (different content hash)
				Name:       "rule2",
				PolicyName: "policy1",
				Action:     model.ActionAccept,
			},
			{
				ID:         "4", // Added
				Name:       "rule3",
				PolicyName: "policy1",
				Action:     model.ActionLog,
			},
		},
	}

	plan := GeneratePlan(current, desired)

	if plan.AddCount != 1 {
		t.Errorf("expected 1 added rule, got %d", plan.AddCount)
	}
	if len(plan.ToAdd) != 1 || plan.ToAdd[0].Name != "rule3" {
		t.Errorf("wrong rule added")
	}

	if plan.RemoveCount != 0 {
		t.Errorf("expected 0 removed rules, got %d", plan.RemoveCount)
	}

	if plan.ModifyCount != 1 {
		t.Errorf("expected 1 modified rule, got %d", plan.ModifyCount)
	}
	if len(plan.ToModify) != 1 || plan.ToModify[0].After.Name != "rule2" {
		t.Errorf("wrong rule modified")
	}
	if len(plan.ToModify[0].Fields) != 1 || plan.ToModify[0].Fields[0] != "action" {
		t.Errorf("wrong field detected as modified: %v", plan.ToModify[0].Fields)
	}

	// Test Remove
	desired2 := &model.CompiledRuleSet{
		Rules: []model.CompiledRule{
			{
				ID:         "1",
				Name:       "rule1",
				PolicyName: "policy1",
				Action:     model.ActionAccept,
			},
		},
	}
	plan2 := GeneratePlan(current, desired2)
	if plan2.RemoveCount != 1 || plan2.ToRemove[0].Name != "rule2" {
		t.Errorf("expected 1 removed rule (rule2), got %d", plan2.RemoveCount)
	}
}
