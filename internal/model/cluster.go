package model

import "time"

// NodeState represents the current state of a Raft node.
type NodeState string

const (
	StateFollower  NodeState = "follower"
	StateCandidate NodeState = "candidate"
	StateLeader    NodeState = "leader"
)

// EntryType represents the type of log entry.
type EntryType string

const (
	EntryPolicyUpdate EntryType = "PolicyUpdate"
	EntryConfigChange EntryType = "ConfigChange"
	EntryNodeJoin     EntryType = "NodeJoin"
	EntryNodeLeave    EntryType = "NodeLeave"
)

// LogEntry is a single entry in the Raft log.
type LogEntry struct {
	Term      uint64    `json:"term"`
	Index     uint64    `json:"index"`
	Type      EntryType `json:"type"`
	Data      []byte    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// NodeStatus represents the current status of a node in the cluster.
type NodeStatus struct {
	ID        string    `json:"id"`
	State     NodeState `json:"state"`
	Backend   string    `json:"backend"`
	Rules     int       `json:"rules"`
	LastSync  time.Time `json:"lastSync"`
	IsHealthy bool      `json:"isHealthy"`
}
