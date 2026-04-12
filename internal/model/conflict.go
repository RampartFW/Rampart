package model

type ConflictType string

const (
	ConflictShadow        ConflictType = "shadow"
	ConflictContradiction ConflictType = "contradiction"
	ConflictRedundancy    ConflictType = "redundancy"
	ConflictSubset        ConflictType = "subset"
	ConflictOverlap       ConflictType = "overlap"
)

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

type Conflict struct {
	Type     ConflictType `json:"type"`
	Severity Severity     `json:"severity"`
	RuleA    CompiledRule `json:"ruleA"`
	RuleB    CompiledRule `json:"ruleB"`
	Message  string       `json:"message"`
}
