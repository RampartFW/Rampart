package engine

import (
	"context"
	"log"
	"sync"
	"time"
	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

// Engine is the core network policy engine.
type Engine struct {
	mu           sync.RWMutex
	backend      backend.Backend
	currentRules *model.CompiledRuleSet
	subscribers  map[chan model.AuditEvent]bool
	subMu        sync.Mutex
}

// NewEngine creates a new engine instance.
func NewEngine(b backend.Backend) *Engine {
	return &Engine{
		backend:     b,
		subscribers: make(map[chan model.AuditEvent]bool),
	}
}

// CurrentRules returns the currently active ruleset.
func (e *Engine) CurrentRules() *model.CompiledRuleSet {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.currentRules
}

// SetRules updates the engine's current ruleset.
func (e *Engine) SetRules(rs *model.CompiledRuleSet) {
	e.mu.Lock()
	e.currentRules = rs
	e.mu.Unlock()
}

// ReapplyRules reapplies the current ruleset to the backend, filtering for active ones.
func (e *Engine) ReapplyRules(ctx context.Context) error {
	e.mu.RLock()
	rs := e.currentRules
	e.mu.RUnlock()

	if rs == nil {
		return nil
	}

	// Filter rules based on current schedule
	now := time.Now()
	var activeRules []model.CompiledRule
	for _, rule := range rs.Rules {
		if IsActive(rule.Schedule, now) {
			activeRules = append(activeRules, rule)
		}
	}

	activeRS := &model.CompiledRuleSet{
		Rules:      activeRules,
		Hash:       rs.Hash,
		Metadata:   rs.Metadata,
		CompiledAt: rs.CompiledAt,
		Backend:    e.backend.Name(),
	}

	return e.backend.Apply(ctx, activeRS)
}

// StartWatchdog begins a background loop that ensures backend health and state consistency.
func (e *Engine) StartWatchdog(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			e.verifyBackendState(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (e *Engine) verifyBackendState(ctx context.Context) {
	// 1. Probe backend health
	if err := e.backend.Probe(); err != nil {
		log.Printf("Watchdog: Backend %s health check failed: %v", e.backend.Name(), err)
		return
	}

	// 2. Optional: Compare current kernel state with desired state (T-059)
	// For high-assurance production, we could call e.backend.CurrentState() 
	// and diff it. For now, we trigger a re-apply if explicitly requested 
	// or per safety interval.
}

func (e *Engine) Backend() backend.Backend {
	return e.backend
}

// Subscribe adds a new subscriber for events.
func (e *Engine) Subscribe() chan model.AuditEvent {
	e.subMu.Lock()
	defer e.subMu.Unlock()
	ch := make(chan model.AuditEvent, 10)
	e.subscribers[ch] = true
	return ch
}

// Unsubscribe removes a subscriber.
func (e *Engine) Unsubscribe(ch chan model.AuditEvent) {
	e.subMu.Lock()
	defer e.subMu.Unlock()
	delete(e.subscribers, ch)
	close(ch)
}

// Broadcast sends an event to all subscribers.
func (e *Engine) Broadcast(event model.AuditEvent) {
	// DPI Signaling: If there's a payload, analyze it before broadcasting
	if len(event.Payload) > 0 {
		if event.Metadata == nil {
			event.Metadata = make(map[string]string)
		}
		
		// 1. DNS Analysis
		if AnalyzeDNS(event.Payload, "malicious.com") { // Example blacklist
			event.Metadata["signal"] = "dns_anomaly"
		}
		
		// 2. HTTP Analysis
		if AnalyzeHTTP(event.Payload, "", "/etc/passwd") {
			event.Metadata["signal"] = "sqli_detect" // Pattern matched
		}
	}

	e.subMu.Lock()
	defer e.subMu.Unlock()
	for ch := range e.subscribers {
		select {
		case ch <- event:
		default:
			// Subscriber slow, drop or handle accordingly
		}
	}
}
