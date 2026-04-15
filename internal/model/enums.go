package model

import (
	"encoding/json"
	"strings"
)

type Direction string

const (
	DirectionInbound  Direction = "inbound"
	DirectionOutbound Direction = "outbound"
	DirectionForward  Direction = "forward"
)

type Action string

const (
	ActionAccept    Action = "accept"
	ActionDrop      Action = "drop"
	ActionReject    Action = "reject"
	ActionLog       Action = "log"
	ActionRateLimit Action = "rate-limit"
)

type Protocol string

const (
	ProtocolTCP    Protocol = "tcp"
	ProtocolUDP    Protocol = "udp"
	ProtocolICMP   Protocol = "icmp"
	ProtocolICMPv6 Protocol = "icmpv6"
	ProtocolAny    Protocol = "any"
)

type ConnState string

const (
	StateNew         ConnState = "new"
	StateEstablished ConnState = "established"
	StateRelated     ConnState = "related"
	StateInvalid     ConnState = "invalid"
)

type IPVersion string

const (
	IPv4   IPVersion = "ipv4"
	IPv6   IPVersion = "ipv6"
	IPBoth IPVersion = "both"
)

func (d Direction) String() string { return string(d) }
func (a Action) String() string    { return string(a) }
func (p Protocol) String() string  { return string(p) }
func (c ConnState) String() string { return string(c) }
func (i IPVersion) String() string { return string(i) }

func (d Direction) MarshalJSON() ([]byte, error) { return json.Marshal(string(d)) }
func (a Action) MarshalJSON() ([]byte, error)    { return json.Marshal(string(a)) }
func (p Protocol) MarshalJSON() ([]byte, error)  { return json.Marshal(string(p)) }
func (c ConnState) MarshalJSON() ([]byte, error) { return json.Marshal(string(c)) }
func (i IPVersion) MarshalJSON() ([]byte, error) { return json.Marshal(string(i)) }

func ActionFromString(s string) Action {
	switch strings.ToLower(s) {
	case "accept":
		return ActionAccept
	case "drop":
		return ActionDrop
	case "reject":
		return ActionReject
	case "log":
		return ActionLog
	case "rate-limit":
		return ActionRateLimit
	default:
		return ActionAccept
	}
}

func ProtocolFromString(s string) Protocol {
	switch strings.ToLower(s) {
	case "tcp":
		return ProtocolTCP
	case "udp":
		return ProtocolUDP
	case "icmp":
		return ProtocolICMP
	case "icmpv6":
		return ProtocolICMPv6
	case "any":
		return ProtocolAny
	default:
		return ProtocolAny
	}
}

func (d *Direction) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*d = Direction(strings.ToLower(s))
	return nil
}

func (a *Action) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*a = Action(strings.ToLower(s))
	return nil
}
