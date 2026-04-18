package engine

import (
	"bytes"
	"strings"

	"github.com/rampartfw/rampart/internal/model"
)

// AnalyzeDNS attempts to parse a DNS query from the packet payload and match it against rules.
func AnalyzeDNS(payload []byte, query string) bool {
	// Simple DNS parser: looking for the domain string in the payload
	// In a real DPI engine, we would parse the DNS header and Question section properly.
	// DNS queries are encoded as [length]label[length]label...
	
	// Convert human-readable query "example.com" to DNS wire format
	parts := strings.Split(query, ".")
	var wireQuery []byte
	for _, p := range parts {
		wireQuery = append(wireQuery, byte(len(p)))
		wireQuery = append(wireQuery, []byte(p)...)
	}
	wireQuery = append(wireQuery, 0) // Null terminator

	return bytes.Contains(payload, wireQuery)
}

// AnalyzeHTTP attempts to match HTTP host or path from the payload.
func AnalyzeHTTP(payload []byte, host string, path string) bool {
	pStr := string(payload)
	
	if host != "" && !strings.Contains(pStr, "Host: "+host) {
		return false
	}
	
	if path != "" && !strings.Contains(pStr, "GET "+path) && !strings.Contains(pStr, "POST "+path) {
		return false
	}
	
	return true
}

// MatchesL7 checks if the packet payload matches the Layer-7 criteria of a rule.
func MatchesL7(rule model.CompiledRule, payload []byte) bool {
	if rule.Match.AppProtocol == "" {
		return true // No L7 criteria
	}

	switch strings.ToLower(rule.Match.AppProtocol) {
	case "dns":
		if rule.Match.DNS != nil && rule.Match.DNS.Query != "" {
			return AnalyzeDNS(payload, rule.Match.DNS.Query)
		}
	case "http":
		if rule.Match.HTTP != nil {
			return AnalyzeHTTP(payload, rule.Match.HTTP.Host, rule.Match.HTTP.Path)
		}
	}

	return false
}
