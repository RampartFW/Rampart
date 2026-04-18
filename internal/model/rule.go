package model

import (
	"net"
	"time"
)

type Rule struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	PolicyID    string            `json:"policyId"`
	Priority    int               `json:"priority"`
	Direction   Direction         `json:"direction"`
	Action      Action            `json:"action"`
	Match       Match             `json:"match"`
	Schedule    *Schedule         `json:"schedule,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	Description string            `json:"description,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	CreatedBy   string            `json:"createdBy"`
	Version     uint64            `json:"version"`
}

type CompiledRule struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	PolicyName  string            `json:"policyName"`
	Priority    int               `json:"priority"`
	Direction   Direction         `json:"direction"`
	Action      Action            `json:"action"`
	Match       CompiledMatch     `json:"match"`
	Log         bool              `json:"log,omitempty"`
	RateLimit   *RateLimit        `json:"rateLimit,omitempty"`
	Schedule    *Schedule         `json:"schedule,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	Description string            `json:"description,omitempty"`
	SourceFile  string            `json:"sourceFile"`
	SourceLine  int               `json:"sourceLine"`
}

type CompiledMatch struct {
	SourceNets  []net.IPNet    `json:"sourceNets,omitempty"`
	DestNets    []net.IPNet    `json:"destNets,omitempty"`
	SourcePorts []PortRange    `json:"sourcePorts,omitempty"`
	DestPorts   []PortRange    `json:"destPorts,omitempty"`
	Protocols   []Protocol     `json:"protocols,omitempty"`
	Interfaces  []string       `json:"interfaces,omitempty"`
	States      []ConnState    `json:"states,omitempty"`
	ICMPTypes   []uint8        `json:"icmpTypes,omitempty"`
	IPVersion    IPVersion      `json:"ipVersion,omitempty"`
	Negated      *CompiledMatch `json:"negated,omitempty"`

	// Layer-7 / DPI Fields (Milestone 21)
	AppProtocol string     `json:"appProtocol,omitempty"` // "http", "tls", "dns"
	HTTP        *HTTPMatch `json:"http,omitempty"`
	TLS         *TLSMatch  `json:"tls,omitempty"`
	DNS         *DNSMatch  `json:"dns,omitempty"`

	// Internal cache for faster conflict detection
	SrcIntervals []IPInterval `json:"-"`
	DstIntervals []IPInterval `json:"-"`
}

type HTTPMatch struct {
	Host   string `json:"host,omitempty"`
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
}

type TLSMatch struct {
	SNI string `json:"sni,omitempty"`
}

type DNSMatch struct {
	Query string `json:"query,omitempty"`
}

type IPInterval struct {
	Start []byte
	End   []byte
}

type Match struct {
	SourceCIDRs []string    `json:"sourceCIDRs,omitempty"`
	DestCIDRs   []string    `json:"destCIDRs,omitempty"`
	SourcePorts []PortRange `json:"sourcePorts,omitempty"`
	DestPorts   []PortRange `json:"destPorts,omitempty"`
	Protocols   []Protocol  `json:"protocols,omitempty"`
	Interfaces  []string    `json:"interfaces,omitempty"`
	States      []ConnState `json:"states,omitempty"`
	Not         *Match      `json:"not,omitempty"`
	ICMPTypes   []uint8     `json:"icmpTypes,omitempty"`
	IPVersion   IPVersion   `json:"ipVersion,omitempty"`
}

type PortRange struct {
	Start uint16 `json:"start"`
	End   uint16 `json:"end"`
}

type RateLimit struct {
	Rate   int    `json:"rate"`
	Per    string `json:"per"`
	Burst  int    `json:"burst"`
	Action Action `json:"action"`
}

type Schedule struct {
	ActiveFrom  *time.Time     `json:"activeFrom,omitempty"`
	ActiveUntil *time.Time     `json:"activeUntil,omitempty"`
	Recurring   *RecurringSpec `json:"recurring,omitempty"`
	WasActive   bool           `json:"-"`
}

type RecurringSpec struct {
	Days      []time.Weekday `json:"days,omitempty"`
	StartTime string         `json:"startTime"`
	EndTime   string         `json:"endTime"`
	Timezone  string         `json:"timezone"`
}

type CompiledRuleSet struct {
	Rules       []CompiledRule `json:"rules"`
	Hash        string         `json:"hash"`
	CompiledAt  time.Time      `json:"compiledAt"`
	SourceFiles []string       `json:"sourceFiles"`
	Metadata    PolicyMetadata `json:"metadata"`
	Backend     string         `json:"backend,omitempty"`
}
