package model

import (
	"net"
	"time"
)

// BackendCapabilities reports what this backend supports
type BackendCapabilities struct {
	IPv4               bool
	IPv6               bool
	RateLimiting       bool
	ConnectionTracking bool
	Logging            bool
	NAT                bool
	PerRuleCounters    bool
	AtomicReplace      bool
	InterfaceFiltering bool
	MarkPackets        bool
	GeoIP              bool
}

type RuleStats struct {
	RuleID  string
	Packets uint64
	Bytes   uint64
	LastHit time.Time
}

type SimulatedPacket struct {
	SourceIP   net.IP
	DestIP     net.IP
	Protocol   Protocol
	SourcePort uint16
	DestPort   uint16
	Direction  Direction
	Interface  string
	State      ConnState
}

type SimulationResult struct {
	Verdict     Action        // Accept, Drop, Reject
	MatchedRule *CompiledRule // nil if default policy
	MatchPath   string        // Human-readable match explanation
	Evaluated   int           // Number of rules evaluated
	Duration    time.Duration // Simulation time
}
