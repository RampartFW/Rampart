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
	scores map[string]int // IP -> Current Risk Score
}

func NewIPSRunner(eng *Engine, as *audit.Store, opts IPSOptions) *IPSRunner {
	return &IPSRunner{
		engine: eng,
		store:  as,
		opts:   opts,
		scores: make(map[string]int),
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
		Statistics: true,
	}

	events, _, err := ips.store.Search(query)
	if err != nil {
		return
	}

	// Dynamic Risk Scoring
	newScores := make(map[string]int)
	for _, e := range events {
		ip := e.Actor.Identity
		if e.Result.Status == "dropped" {
			newScores[ip] += 1 // Basic drop adds 1 point
		}
		
		// Check for signals in Metadata
		if sig, ok := e.Metadata["signal"]; ok {
			switch sig {
			case "dns_flood":   newScores[ip] += 10
			case "l7_violation": newScores[ip] += 20
			case "sqli_detect":  newScores[ip] += 50
			}
		}
	}

	for ip, score := range newScores {
		if score >= ips.opts.Threshold {
			ips.triggerBlock(ip, score)
		}
	}
	ips.scores = newScores
}

func (ips *IPSRunner) triggerBlock(ip string, score int) {
	log.Printf("IPS: ALERT! Threat score for %s reached %d. Locking down node.", ip, score)
	
	ips.engine.Broadcast(model.AuditEvent{
		ID:        model.GenerateUUIDv7(),
		Action:    model.AuditApply,
		Timestamp: time.Now(),
		Actor:     model.AuditActor{Type: "ips", Identity: "autonomous-sentinel"},
		Resource:  model.AuditResource{Type: "ip-ban", Name: ip},
		Metadata: map[string]string{
			"reason": "risk_score_exceeded",
			"score":  fmt.Sprintf("%d", score),
		},
		Result: model.AuditResultSuccess,
	})
}
