package model

import "time"

type Snapshot struct {
	ID          string            `json:"id"`
	CreatedAt   time.Time         `json:"createdAt"`
	CreatedBy   string            `json:"createdBy"`
	Trigger     string            `json:"trigger"`
	Description string            `json:"description"`
	PolicyHash  string            `json:"policyHash"`
	RuleCount   int               `json:"ruleCount"`
	Backend     string            `json:"backend"`
	State       []byte            `json:"-"`
	Size        int64             `json:"size"`
	Metadata    map[string]string `json:"metadata"`
	Filename    string            `json:"filename"`
}
