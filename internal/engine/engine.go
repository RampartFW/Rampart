package engine

import (
	"context"
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

	// Create a new temporary ruleset for applying
	activeRS := &model.CompiledRuleSet{
		Rules:      activeRules,
		Hash:       rs.Hash, // Keep original hash for tracking
		Metadata:   rs.Metadata,
		CompiledAt: rs.CompiledAt,
	}

	return e.backend.Apply(ctx, activeRS)
}

// Backend returns the engine's backend.
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
