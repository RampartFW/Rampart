package engine

import (
	"net"
	"testing"

	"github.com/rampartfw/rampart/internal/model"
)

func TestDetectConflicts(t *testing.T) {
	_, anyIP, _ := net.ParseCIDR("0.0.0.0/0")
	_, subnetA, _ := net.ParseCIDR("10.0.0.0/24")
	_, subnetB, _ := net.ParseCIDR("10.0.0.0/8")

	tests := []struct {
		name          string
		rules         []model.CompiledRule
		expectedCount int
		expectedType  model.ConflictType
	}{
		{
			name: "Redundancy - Identical rules",
			rules: []model.CompiledRule{
				{
					Name:      "rule1",
					Priority:  10,
					Direction: model.DirectionInbound,
					Action:    model.ActionAccept,
					Match: model.CompiledMatch{
						Protocols:  []model.Protocol{model.ProtocolTCP},
						DestPorts:  []model.PortRange{{Start: 80, End: 80}},
						SourceNets: []net.IPNet{*anyIP},
					},
				},
				{
					Name:      "rule2",
					Priority:  10,
					Direction: model.DirectionInbound,
					Action:    model.ActionAccept,
					Match: model.CompiledMatch{
						Protocols:  []model.Protocol{model.ProtocolTCP},
						DestPorts:  []model.PortRange{{Start: 80, End: 80}},
						SourceNets: []net.IPNet{*anyIP},
					},
				},
			},
			expectedCount: 1,
			expectedType:  model.ConflictRedundancy,
		},
		{
			name: "Shadowing - Higher priority rule covers lower",
			rules: []model.CompiledRule{
				{
					Name:      "rule1",
					Priority:  10,
					Direction: model.DirectionInbound,
					Action:    model.ActionAccept,
					Match: model.CompiledMatch{
						Protocols:  []model.Protocol{model.ProtocolTCP},
						DestPorts:  []model.PortRange{{Start: 80, End: 80}},
						SourceNets: []net.IPNet{*anyIP},
					},
				},
				{
					Name:      "rule2",
					Priority:  20,
					Direction: model.DirectionInbound,
					Action:    model.ActionDrop,
					Match: model.CompiledMatch{
						Protocols:  []model.Protocol{model.ProtocolTCP},
						DestPorts:  []model.PortRange{{Start: 80, End: 80}},
						SourceNets: []net.IPNet{*subnetA},
					},
				},
			},
			expectedCount: 1,
			expectedType:  model.ConflictShadow,
		},
		{
			name: "Contradiction - Same priority, overlap, different action",
			rules: []model.CompiledRule{
				{
					Name:      "rule1",
					Priority:  10,
					Direction: model.DirectionInbound,
					Action:    model.ActionAccept,
					Match: model.CompiledMatch{
						Protocols:  []model.Protocol{model.ProtocolTCP},
						DestPorts:  []model.PortRange{{Start: 80, End: 80}},
						SourceNets: []net.IPNet{*subnetA},
					},
				},
				{
					Name:      "rule2",
					Priority:  10,
					Direction: model.DirectionInbound,
					Action:    model.ActionDrop,
					Match: model.CompiledMatch{
						Protocols:  []model.Protocol{model.ProtocolTCP},
						DestPorts:  []model.PortRange{{Start: 80, End: 80}},
						SourceNets: []net.IPNet{*subnetA},
					},
				},
			},
			expectedCount: 1,
			expectedType:  model.ConflictContradiction,
		},
		{
			name: "Subset - Same action, small covers large",
			rules: []model.CompiledRule{
				{
					Name:      "rule1",
					Priority:  10,
					Direction: model.DirectionInbound,
					Action:    model.ActionAccept,
					Match: model.CompiledMatch{
						Protocols:  []model.Protocol{model.ProtocolTCP},
						DestPorts:  []model.PortRange{{Start: 0, End: 65535}},
						SourceNets: []net.IPNet{*subnetB},
					},
				},
				{
					Name:      "rule2",
					Priority:  20,
					Direction: model.DirectionInbound,
					Action:    model.ActionAccept,
					Match: model.CompiledMatch{
						Protocols:  []model.Protocol{model.ProtocolTCP},
						DestPorts:  []model.PortRange{{Start: 80, End: 80}},
						SourceNets: []net.IPNet{*subnetA},
					},
				},
			},
			expectedCount: 1,
			expectedType:  model.ConflictSubset,
		},
		{
			name: "Overlap - Partial overlap, different action",
			rules: []model.CompiledRule{
				{
					Name:      "rule1",
					Priority:  10,
					Direction: model.DirectionInbound,
					Action:    model.ActionAccept,
					Match: model.CompiledMatch{
						Protocols:  []model.Protocol{model.ProtocolTCP},
						DestPorts:  []model.PortRange{{Start: 80, End: 80}},
						SourceNets: []net.IPNet{*subnetB},
					},
				},
				{
					Name:      "rule2",
					Priority:  20,
					Direction: model.DirectionInbound,
					Action:    model.ActionDrop,
					Match: model.CompiledMatch{
						Protocols:  []model.Protocol{model.ProtocolTCP},
						DestPorts:  []model.PortRange{{Start: 70, End: 90}},
						SourceNets: []net.IPNet{*subnetA},
					},
				},
			},
			expectedCount: 1,
			expectedType:  model.ConflictOverlap,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conflicts := DetectConflicts(tt.rules)
			if len(conflicts) != tt.expectedCount {
				t.Errorf("Expected %d conflicts, got %d", tt.expectedCount, len(conflicts))
			}
			if len(conflicts) > 0 && conflicts[0].Type != tt.expectedType {
				t.Errorf("Expected conflict type %s, got %s", tt.expectedType, conflicts[0].Type)
			}
		})
	}
}

func TestCidrsSubset(t *testing.T) {
	_, anyIP, _ := net.ParseCIDR("0.0.0.0/0")
	_, subnet10, _ := net.ParseCIDR("10.0.0.0/8")
	_, subnet10_0_0, _ := net.ParseCIDR("10.0.0.0/24")
	_, subnet10_0_0_1, _ := net.ParseCIDR("10.0.0.1/32")
	_, subnet192, _ := net.ParseCIDR("192.168.1.0/24")

	tests := []struct {
		name   string
		small  []net.IPNet
		large  []net.IPNet
		expect bool
	}{
		{"Any is any", []net.IPNet{*anyIP}, []net.IPNet{*anyIP}, true},
		{"Subnet in any", []net.IPNet{*subnet10}, []net.IPNet{*anyIP}, true},
		{"Any in subnet", []net.IPNet{*anyIP}, []net.IPNet{*subnet10}, false},
		{"Subnet in larger subnet", []net.IPNet{*subnet10_0_0}, []net.IPNet{*subnet10}, true},
		{"IP in subnet", []net.IPNet{*subnet10_0_0_1}, []net.IPNet{*subnet10_0_0}, true},
		{"Unrelated subnets", []net.IPNet{*subnet192}, []net.IPNet{*subnet10}, false},
		{"Multiple small in one large", []net.IPNet{*subnet10_0_0, *subnet10_0_0_1}, []net.IPNet{*subnet10}, true},
		{"Small not in any of large", []net.IPNet{*subnet192}, []net.IPNet{*subnet10, *subnet10_0_0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cidrsSubset(tt.small, tt.large)
			if result != tt.expect {
				t.Errorf("cidrsSubset() = %v, expect %v", result, tt.expect)
			}
		})
	}
}
