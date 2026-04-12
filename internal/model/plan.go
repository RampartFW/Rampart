package model

type ExecutionPlan struct {
	ToAdd    []CompiledRule     `json:"toAdd"`
	ToRemove []CompiledRule     `json:"toRemove"`
	ToModify []RuleModification `json:"toModify"`
	Warnings []Conflict         `json:"warnings"`
	Errors   []Conflict         `json:"errors"`

	CurrentRuleCount int `json:"currentRuleCount"`
	PlannedRuleCount int `json:"plannedRuleCount"`
	AddCount         int `json:"addCount"`
	RemoveCount      int `json:"removeCount"`
	ModifyCount      int `json:"modifyCount"`
}

type RuleModification struct {
	Before CompiledRule `json:"before"`
	After  CompiledRule `json:"after"`
	Fields []string     `json:"fields"`
}
