package model

import (
	"encoding/json"
	"time"
)

type AuditAction string

const (
	AuditApply    AuditAction = "policy.apply"
	AuditRollback AuditAction = "policy.rollback"
	AuditPlan     AuditAction = "policy.plan"
	AuditSimulate AuditAction = "policy.simulate"
	AuditSnapshot AuditAction = "snapshot.create"
	AuditFlush    AuditAction = "policy.flush"
	AuditSync     AuditAction = "cluster.sync"
	AuditJoin     AuditAction = "cluster.join"
	AuditLeave    AuditAction = "cluster.leave"
)

type AuditEvent struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	NodeID    string            `json:"nodeId"`
	Actor     AuditActor        `json:"actor"`
	Action    AuditAction       `json:"action"`
	Resource  AuditResource     `json:"resource"`
	Before    json.RawMessage   `json:"before,omitempty"`
	After     json.RawMessage   `json:"after,omitempty"`
	Result    AuditResult       `json:"result"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Payload   []byte            `json:"payload,omitempty"` // For DPI analysis
	ChainHash string            `json:"chainHash,omitempty"`
}

type AuditActor struct {
	Type     string `json:"type"`     // "user", "api", "system", "mcp", "raft-sync"
	Identity string `json:"identity"` // Username, API key ID, "system:scheduler"
	SourceIP string `json:"sourceIp,omitempty"`
}

type AuditResource struct {
	Type string `json:"type"` // "policy", "rule", "snapshot", "cluster"
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type AuditResult struct {
	Status  string `json:"status"` // "success", "failure", "dry-run"
	Message string `json:"message,omitempty"`
}

var (
	AuditResultSuccess = AuditResult{Status: "success"}
	AuditResultFailure = AuditResult{Status: "failure"}
	AuditResultDryRun  = AuditResult{Status: "dry-run"}
)
