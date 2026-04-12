package engine

import (
	"bytes"
	"fmt"
	"math/big"
	"net"
	"reflect"

	"github.com/rampartfw/rampart/internal/model"
)

// DetectConflicts performs pairwise conflict detection for a set of compiled rules.
func DetectConflicts(rules []model.CompiledRule) []model.Conflict {
	var conflicts []model.Conflict

	for i := 0; i < len(rules); i++ {
		for j := i + 1; j < len(rules); j++ {
			if c := checkPair(rules[i], rules[j]); c != nil {
				conflicts = append(conflicts, *c)
			}
		}
	}

	return conflicts
}

func checkPair(a, b model.CompiledRule) *model.Conflict {
	// Same direction required for conflict
	if a.Direction != b.Direction {
		return nil
	}

	// Check protocol overlap
	if !protocolsOverlap(a.Match.Protocols, b.Match.Protocols) {
		return nil
	}

	// Check port overlap
	if !portsOverlap(a.Match.SourcePorts, b.Match.SourcePorts) {
		return nil
	}
	if !portsOverlap(a.Match.DestPorts, b.Match.DestPorts) {
		return nil
	}

	// Check CIDR overlap
	if !cidrsOverlapCached(a.Match.SrcIntervals, b.Match.SrcIntervals) {
		return nil
	}
	if !cidrsOverlapCached(a.Match.DstIntervals, b.Match.DstIntervals) {
		return nil
	}

	// Check interface overlap
	if !interfacesOverlap(a.Match.Interfaces, b.Match.Interfaces) {
		return nil
	}

	// Check state overlap
	if !statesOverlap(a.Match.States, b.Match.States) {
		return nil
	}

	// Overlap found — classify the conflict
	if a.Action == b.Action {
		if matchEquals(a.Match, b.Match) {
			return &model.Conflict{
				Type:     model.ConflictRedundancy,
				Severity: model.SeverityInfo,
				RuleA:    a,
				RuleB:    b,
				Message:  fmt.Sprintf("Rule %q and %q are identical.", a.Name, b.Name),
			}
		}
		if matchSubset(b.Match, a.Match) {
			return &model.Conflict{
				Type:     model.ConflictSubset,
				Severity: model.SeverityInfo,
				RuleA:    a,
				RuleB:    b,
				Message:  fmt.Sprintf("Rule %q is a subset of %q and has the same action.", b.Name, a.Name),
			}
		}
		// If they overlap but aren't subsets/equals, we could still report it, 
		// but redundancy/subset are usually what's interesting for same-action.
		return nil
	}

	if a.Priority == b.Priority {
		return &model.Conflict{
			Type:     model.ConflictContradiction,
			Severity: model.SeverityError,
			RuleA:    a,
			RuleB:    b,
			Message:  fmt.Sprintf("Rules %q and %q have the same priority but different actions.", a.Name, b.Name),
		}
	}

	// Different priority, different action
	// If a higher priority rule (a) covers a lower priority rule (b), it's a shadow
	if matchSubset(b.Match, a.Match) {
		return &model.Conflict{
			Type:     model.ConflictShadow,
			Severity: model.SeverityWarning,
			RuleA:    a,
			RuleB:    b,
			Message:  fmt.Sprintf("Rule %q is shadowed by higher priority rule %q.", b.Name, a.Name),
		}
	}

	return &model.Conflict{
		Type:     model.ConflictOverlap,
		Severity: model.SeverityWarning,
		RuleA:    a,
		RuleB:    b,
		Message:  fmt.Sprintf("Rules %q and %q overlap with different actions.", a.Name, b.Name),
	}
}

func protocolsOverlap(a, b []model.Protocol) bool {
	if len(a) == 0 || len(b) == 0 {
		return true
	}
	for _, p1 := range a {
		if p1 == model.ProtocolAny {
			return true
		}
		for _, p2 := range b {
			if p2 == model.ProtocolAny || p1 == p2 {
				return true
			}
		}
	}
	return false
}

func portsOverlap(a, b []model.PortRange) bool {
	if len(a) == 0 || len(b) == 0 {
		return true
	}
	for _, r1 := range a {
		for _, r2 := range b {
			if r1.Start <= r2.End && r2.Start <= r1.End {
				return true
			}
		}
	}
	return false
}

func interfacesOverlap(a, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return true
	}
	for _, i1 := range a {
		for _, i2 := range b {
			if i1 == i2 {
				return true
			}
		}
	}
	return false
}

func statesOverlap(a, b []model.ConnState) bool {
	if len(a) == 0 || len(b) == 0 {
		return true
	}
	for _, s1 := range a {
		for _, s2 := range b {
			if s1 == s2 {
				return true
			}
		}
	}
	return false
}

func cidrsOverlap(a, b []net.IPNet) bool {
	if len(a) == 0 || len(b) == 0 {
		return true
	}
	for _, n1 := range a {
		for _, n2 := range b {
			start1, end1 := cidrToInterval(n1)
			start2, end2 := cidrToInterval(n2)
			if start1.Cmp(end2) <= 0 && start2.Cmp(end1) <= 0 {
				return true
			}
		}
	}
	return false
}

func cidrsOverlapCached(a, b []model.IPInterval) bool {
	if len(a) == 0 || len(b) == 0 {
		return true
	}
	for _, i1 := range a {
		for _, i2 := range b {
			if bytes.Compare(i1.Start, i2.End) <= 0 && bytes.Compare(i2.Start, i1.End) <= 0 {
				return true
			}
		}
	}
	return false
}

// Convert CIDR to interval
func cidrToInterval(ipNet net.IPNet) (*big.Int, *big.Int) {
	start := ipToInt(ipNet.IP)
	ones, bits := ipNet.Mask.Size()
	size := new(big.Int).Lsh(big.NewInt(1), uint(bits-ones))
	end := new(big.Int).Add(start, size)
	end.Sub(end, big.NewInt(1))
	return start, end
}

func matchEquals(a, b model.CompiledMatch) bool {
	return reflect.DeepEqual(a, b)
}

func matchSubset(small, large model.CompiledMatch) bool {
	// small is subset of large if all conditions in small are more restrictive or equal to large.
	
	// Protocols
	if len(large.Protocols) > 0 {
		if len(small.Protocols) == 0 {
			// small matches any protocol, but large only matches some. Not a subset.
			// (Unless large also matches any, but we checked len(large.Protocols) > 0)
			// Wait, in our model, empty means ANY.
			return false 
		}
		// Every protocol in small must be in large
		for _, pS := range small.Protocols {
			found := false
			for _, pL := range large.Protocols {
				if pL == model.ProtocolAny || pS == pL {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	// CIDRs
	if !cidrsSubset(small.SourceNets, large.SourceNets) {
		return false
	}
	if !cidrsSubset(small.DestNets, large.DestNets) {
		return false
	}

	// Ports
	if !portsSubset(small.SourcePorts, large.SourcePorts) {
		return false
	}
	if !portsSubset(small.DestPorts, large.DestPorts) {
		return false
	}

	// Interfaces
	if len(large.Interfaces) > 0 {
		if len(small.Interfaces) == 0 {
			return false
		}
		for _, iS := range small.Interfaces {
			if !contains(large.Interfaces, iS) {
				return false
			}
		}
	}

	// States
	if len(large.States) > 0 {
		if len(small.States) == 0 {
			return false
		}
		for _, sS := range small.States {
			if !containsState(large.States, sS) {
				return false
			}
		}
	}

	return true
}

func cidrsSubset(small, large []net.IPNet) bool {
	if len(large) == 0 {
		return true // large is ANY
	}
	if len(small) == 0 {
		return false // small is ANY, large is not
	}
	// Every net in small must be contained in at least one net in large
	for _, nS := range small {
		found := false
		for _, nL := range large {
			if nL.Contains(nS.IP) {
				// We also need to check if the mask is at least as restrictive
				onesS, _ := nS.Mask.Size()
				onesL, _ := nL.Mask.Size()
				if onesS >= onesL {
					found = true
					break
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func portsSubset(small, large []model.PortRange) bool {
	if len(large) == 0 {
		return true // large is ANY
	}
	if len(small) == 0 {
		return false // small is ANY, large is not
	}
	// Every range in small must be fully contained within some range in large
	for _, rS := range small {
		found := false
		for _, rL := range large {
			if rS.Start >= rL.Start && rS.End <= rL.End {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func containsState(slice []model.ConnState, s model.ConnState) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func ipToInt(ip net.IP) *big.Int {
	if ip.To4() != nil {
		return big.NewInt(0).SetBytes(ip.To4())
	}
	return big.NewInt(0).SetBytes(ip.To16())
}
