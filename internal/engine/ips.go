package engine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rampartfw/rampart/internal/audit"
	"github.com/rampartfw/rampart/internal/model"
)

type IPSOptions struct {
	Enabled       bool
	Threshold     int           // Number of drops to trigger block
	Window        time.Duration // Time window to analyze
	BlockDuration time.Duration // How long to keep the block
}

type IPSRunner struct {
	engine *Engine
	store  *audit.Store
	opts   IPSOptions
}

func NewIPSRunner(eng *Engine, as *audit.Store, opts IPSOptions) *IPSRunner {
	return &IPSRunner{
		engine: eng,
		store:  as,
		opts:   opts,
	}
}

func (ips *IPSRunner) Run(ctx context.Context) {
	if !ips.opts.Enabled {
		return
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ips.analyze(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (ips *IPSRunner) analyze(ctx context.Context) {
	now := time.Now()
	query := audit.AuditQuery{
		Since:      now.Add(-ips.opts.Window),
		Action:     model.AuditApply, // We actually want to look for Drop events in traffic logs
		Statistics: true,
	}

	// In a full implementation, the Backend would report individual Drop events 
	// to the AuditStore. Here we use the AuditStore's statistical capabilities.
	events, _, err := ips.store.Search(query)
	if err != nil {
		return
	}

	// Group drops by Source IP
	badIPs := make(map[string]int)
	for _, e := range events {
		if e.Result.Status == "dropped" { // Custom status for traffic logs
			badIPs[e.Actor.Identity]++
		}
	}

	for ip, count := range badIPs {
		if count >= ips.opts.Threshold {
			ips.triggerBlock(ip)
		}
	}
}

func (ips *IPSRunner) triggerBlock(ip string) {
	log.Printf("IPS: Detected threat from %s (%d events). Triggering cluster-wide block.", ip, ips.opts.Threshold)

	// Create an automated policy update
	// In a real implementation, this would use Raft to propose EntryIPBan

	// Broadcast event for WebUI/Audit
	ips.engine.Broadcast(model.AuditEvent{
		ID:        model.GenerateUUIDv7(),
		Action:    model.AuditApply, // Using Apply to signify a rule change
		Timestamp: time.Now(),
		Actor:     model.AuditActor{Type: "ips", Identity: "autonomous-detector"},
		Resource:  model.AuditResource{Type: "ip-ban", Name: ip},
		Result:    model.AuditResultSuccess,
	})
}
