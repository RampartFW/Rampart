package engine

import (
	"net"
	"testing"

	"github.com/rampartfw/rampart/internal/model"
)

func TestSimulator(t *testing.T) {
	_, ipNet10, _ := net.ParseCIDR("10.0.0.0/8")
	_, ipNet192, _ := net.ParseCIDR("192.168.1.0/24")

	rules := []model.CompiledRule{
		{
			ID:        "1",
			Name:      "Allow SSH from office",
			Priority:  10,
			Direction: model.DirectionInbound,
			Action:    model.ActionAccept,
			Match: model.CompiledMatch{
				Protocols:  []model.Protocol{model.ProtocolTCP},
				DestPorts:  []model.PortRange{{Start: 22, End: 22}},
				SourceNets: []net.IPNet{*ipNet10},
			},
		},
		{
			ID:        "2",
			Name:      "Block all to internal DB",
			Priority:  20,
			Direction: model.DirectionInbound,
			Action:    model.ActionDrop,
			Match: model.CompiledMatch{
				DestNets: []net.IPNet{*ipNet192},
			},
		},
	}

	t.Run("Match_First_Rule", func(t *testing.T) {
		pkt := model.SimulatedPacket{
			SourceIP:  net.ParseIP("10.0.0.5"),
			DestIP:    net.ParseIP("172.16.0.1"),
			Protocol:  model.ProtocolTCP,
			DestPort:  22,
			Direction: model.DirectionInbound,
		}
		res := Simulate(rules, pkt)
		if res.Verdict != model.ActionAccept {
			t.Errorf("expected ACCEPT, got %v", res.Verdict)
		}
		if res.MatchedRule.ID != "1" {
			t.Errorf("expected rule 1, got %v", res.MatchedRule.ID)
		}
	})

	t.Run("Match_Second_Rule", func(t *testing.T) {
		pkt := model.SimulatedPacket{
			SourceIP:  net.ParseIP("8.8.8.8"),
			DestIP:    net.ParseIP("192.168.1.100"),
			Protocol:  model.ProtocolTCP,
			DestPort:  3306,
			Direction: model.DirectionInbound,
		}
		res := Simulate(rules, pkt)
		if res.Verdict != model.ActionDrop {
			t.Errorf("expected DROP, got %v", res.Verdict)
		}
		if res.MatchedRule.ID != "2" {
			t.Errorf("expected rule 2, got %v", res.MatchedRule.ID)
		}
	})

	t.Run("Default_Deny", func(t *testing.T) {
		pkt := model.SimulatedPacket{
			SourceIP:  net.ParseIP("8.8.8.8"),
			DestIP:    net.ParseIP("1.1.1.1"),
			Protocol:  model.ProtocolUDP,
			DestPort:  53,
			Direction: model.DirectionInbound,
		}
		res := Simulate(rules, pkt)
		if res.Verdict != model.ActionDrop {
			t.Errorf("expected default DROP, got %v", res.Verdict)
		}
		if res.MatchedRule != nil {
			t.Errorf("expected no matched rule, got %v", res.MatchedRule.ID)
		}
	})

	t.Run("Wrong_Direction", func(t *testing.T) {
		pkt := model.SimulatedPacket{
			SourceIP:  net.ParseIP("10.0.0.5"),
			DestIP:    net.ParseIP("172.16.0.1"),
			Protocol:  model.ProtocolTCP,
			DestPort:  22,
			Direction: model.DirectionOutbound,
		}
		res := Simulate(rules, pkt)
		if res.Verdict != model.ActionDrop {
			t.Errorf("expected DROP for outbound, got %v", res.Verdict)
		}
	})
}
