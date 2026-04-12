package nftables

import (
	"net"
	"strings"
	"testing"

	"github.com/rampartfw/rampart/internal/model"
)

func TestGenerateScript(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("10.0.1.0/24")
	
	rules := []model.CompiledRule{
		{
			Name:      "allow-ssh",
			Direction: model.DirectionInbound,
			Action:    model.ActionAccept,
			Match: model.CompiledMatch{
				Protocols: []model.Protocol{model.ProtocolTCP},
				DestPorts: []model.PortRange{{Start: 22, End: 22}},
				SourceNets: []net.IPNet{*ipnet},
			},
		},
		{
			Name:      "deny-web",
			Direction: model.DirectionInbound,
			Action:    model.ActionDrop,
			Match: model.CompiledMatch{
				Protocols: []model.Protocol{model.ProtocolTCP},
				DestPorts: []model.PortRange{{Start: 80, End: 80}, {Start: 443, End: 443}},
			},
			Log: true,
		},
	}

	rs := &model.CompiledRuleSet{
		Rules: rules,
	}

	script := generateScript(rs)

	// Basic checks
	if !strings.Contains(script, "table inet rampart") {
		t.Errorf("script missing table definition")
	}
	if !strings.Contains(script, "chain input") {
		t.Errorf("script missing input chain")
	}
	if !strings.Contains(script, "allow-ssh") {
		t.Errorf("script missing allow-ssh rule")
	}
	if !strings.Contains(script, "ip saddr 10.0.1.0/24") {
		t.Errorf("script missing source CIDR for allow-ssh")
	}
	if !strings.Contains(script, "tcp dport 22") {
		t.Errorf("script missing dport 22 for allow-ssh")
	}
	if !strings.Contains(script, "log prefix \"rampart:deny-web: \"") {
		t.Errorf("script missing log prefix for deny-web")
	}
}

func TestRenderNets(t *testing.T) {
	_, v4, _ := net.ParseCIDR("1.2.3.4/32")
	_, v6, _ := net.ParseCIDR("2001:db8::/32")

	res := renderNets("saddr", []net.IPNet{*v4})
	if res != "ip saddr 1.2.3.4/32" {
		t.Errorf("expected ip saddr 1.2.3.4/32, got %s", res)
	}

	res = renderNets("daddr", []net.IPNet{*v6})
	if res != "ip6 daddr 2001:db8::/32" {
		t.Errorf("expected ip6 daddr 2001:db8::/32, got %s", res)
	}

	res = renderNets("saddr", []net.IPNet{*v4, *v6})
	if !strings.Contains(res, "ip saddr 1.2.3.4/32") || !strings.Contains(res, "ip6 saddr 2001:db8::/32") {
		t.Errorf("expected both v4 and v6 in result, got %s", res)
	}
}
