package snapshot

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/rampartfw/rampart/internal/model"
)

func generateRules(count int) *model.CompiledRuleSet {
	rules := make([]model.CompiledRule, count)
	for i := 0; i < count; i++ {
		rules[i] = model.CompiledRule{
			ID:         fmt.Sprintf("rule-%d", i),
			Name:       fmt.Sprintf("Allow HTTP %d", i),
			Priority:   i,
			Direction:  model.DirectionInbound,
			Action:     model.ActionAccept,
			Match: model.CompiledMatch{
				Protocols: []model.Protocol{model.ProtocolTCP},
				DestPorts: []model.PortRange{{Start: 80, End: 80}},
				SourceNets: []net.IPNet{
					{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
				},
			},
		}
	}
	return &model.CompiledRuleSet{
		Rules: rules,
		Hash:  "fake-hash",
	}
}

func BenchmarkSnapshotCreate(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "rampart-bench-*")
	defer os.RemoveAll(tmpDir)

	store, _ := NewStore(tmpDir)
	rules := generateRules(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := store.Create("bench", "benchmark", rules)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSnapshotRestore(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "rampart-bench-*")
	defer os.RemoveAll(tmpDir)

	store, _ := NewStore(tmpDir)
	rules := generateRules(1000)
	snap, _ := store.Create("bench", "benchmark", rules)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := store.Load(snap.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}
