package engine

import (
	"reflect"

	"github.com/rampartfw/rampart/internal/model"
)

// GeneratePlan compares current and desired CompiledRuleSets and produces an ExecutionPlan.
func GeneratePlan(current, desired *model.CompiledRuleSet) *model.ExecutionPlan {
	plan := &model.ExecutionPlan{
		CurrentRuleCount: len(current.Rules),
		PlannedRuleCount: len(desired.Rules),
	}

	currentMap := make(map[string]model.CompiledRule)
	for _, r := range current.Rules {
		key := r.PolicyName + ":" + r.Name
		currentMap[key] = r
	}

	desiredMap := make(map[string]model.CompiledRule)
	for _, r := range desired.Rules {
		key := r.PolicyName + ":" + r.Name
		desiredMap[key] = r
	}

	// Detect Add and Modify
	for key, desiredRule := range desiredMap {
		if currentRule, exists := currentMap[key]; exists {
			if currentRule.ID != desiredRule.ID {
				plan.ToModify = append(plan.ToModify, model.RuleModification{
					Before: currentRule,
					After:  desiredRule,
					Fields: diffFields(currentRule, desiredRule),
				})
			}
		} else {
			plan.ToAdd = append(plan.ToAdd, desiredRule)
		}
	}

	// Detect Remove
	for key, currentRule := range currentMap {
		if _, exists := desiredMap[key]; !exists {
			plan.ToRemove = append(plan.ToRemove, currentRule)
		}
	}

	plan.AddCount = len(plan.ToAdd)
	plan.RemoveCount = len(plan.ToRemove)
	plan.ModifyCount = len(plan.ToModify)

	return plan
}

func diffFields(a, b model.CompiledRule) []string {
	var fields []string

	if a.Priority != b.Priority {
		fields = append(fields, "priority")
	}
	if a.Direction != b.Direction {
		fields = append(fields, "direction")
	}
	if a.Action != b.Action {
		fields = append(fields, "action")
	}
	if a.Log != b.Log {
		fields = append(fields, "log")
	}
	if !reflect.DeepEqual(a.Match, b.Match) {
		fields = append(fields, "match")
	}
	if !reflect.DeepEqual(a.RateLimit, b.RateLimit) {
		fields = append(fields, "rateLimit")
	}
	if !reflect.DeepEqual(a.Schedule, b.Schedule) {
		fields = append(fields, "schedule")
	}
	if !reflect.DeepEqual(a.Tags, b.Tags) {
		fields = append(fields, "tags")
	}
	if a.Description != b.Description {
		fields = append(fields, "description")
	}

	return fields
}
