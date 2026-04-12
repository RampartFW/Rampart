package engine

import (
	"fmt"
	"net"
	"testing"

	"github.com/rampartfw/rampart/internal/model"
)

func generateRules(count int) []model.CompiledRule {
	rules := make([]model.CompiledRule, count)
	for i := 0; i < count; i++ {
		ip := net.ParseIP(fmt.Sprintf("10.0.%d.%d", (i/256)%256, i%256))
		_, ipNet, _ := net.ParseCIDR(ip.String() + "/32")
		
		start, end := cidrToInterval(*ipNet)

		rules[i] = model.CompiledRule{
			ID:         fmt.Sprintf("rule-%d", i),
			Name:       fmt.Sprintf("Allow HTTP %d", i),
			Priority:   i,
			Direction:  model.DirectionInbound,
			Action:     model.ActionAccept,
			Match: model.CompiledMatch{
				Protocols: []model.Protocol{model.ProtocolTCP},
				DestPorts: []model.PortRange{{Start: uint16(i % 65535), End: uint16(i % 65535)}},
				SourceNets: []net.IPNet{*ipNet},
				SrcIntervals: []model.IPInterval{
					{Start: start.Bytes(), End: end.Bytes()},
				},
			},
		}
	}
	return rules
}

func BenchmarkCompile100Rules(b *testing.B) {
	ps := &model.PolicySetYAML{
		Metadata: model.PolicyMetadata{Name: "benchmark"},
		Policies: []model.PolicyYAML{
			{
				Name:     "p1",
				Priority: 100,
				Rules:    make([]model.RuleYAML, 100),
			},
		},
	}
	for i := 0; i < 100; i++ {
		ps.Policies[0].Rules[i] = model.RuleYAML{
			Name:   fmt.Sprintf("rule-%d", i),
			Action: model.ActionAccept,
			Match: model.MatchYAML{
				Protocol:    "tcp",
				DestPorts:   80,
				SourceCIDRs: []string{"10.0.0.0/8"},
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Compile(ps, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompile10000Rules(b *testing.B) {
	ps := &model.PolicySetYAML{
		Metadata: model.PolicyMetadata{Name: "benchmark"},
		Policies: []model.PolicyYAML{
			{
				Name:     "p1",
				Priority: 100,
				Rules:    make([]model.RuleYAML, 10000),
			},
		},
	}
	for i := 0; i < 10000; i++ {
		ps.Policies[0].Rules[i] = model.RuleYAML{
			Name:   fmt.Sprintf("rule-%d", i),
			Action: model.ActionAccept,
			Match: model.MatchYAML{
				Protocol:    "tcp",
				DestPorts:   80,
				SourceCIDRs: []string{"10.0.0.0/8"},
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Compile(ps, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConflictDetection1000Rules(b *testing.B) {
	rules := generateRules(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectConflicts(rules)
	}
}

func BenchmarkSimulatePacket(b *testing.B) {
	rules := generateRules(100)
	pkt := model.SimulatedPacket{
		SourceIP:   net.ParseIP("10.1.2.3"),
		DestIP:     net.ParseIP("192.168.1.1"),
		Protocol:   model.ProtocolTCP,
		DestPort:   80,
		Direction:  model.DirectionInbound,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Simulate(rules, pkt)
	}
}
