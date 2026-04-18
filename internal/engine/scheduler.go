package engine

import (
	"context"
	"log"
	"time"

	"github.com/rampartfw/rampart/internal/model"
)

// IsActive checks if a schedule is currently active at the given time.
func IsActive(sched *model.Schedule, now time.Time) bool {
	if sched == nil {
		return true
	}

	// One-time schedule
	if sched.ActiveFrom != nil && now.Before(*sched.ActiveFrom) {
		return false
	}
	if sched.ActiveUntil != nil && now.After(*sched.ActiveUntil) {
		return false
	}

	// Recurring schedule
	if sched.Recurring != nil {
		var loc *time.Location
		var err error
		if sched.Recurring.Timezone != "" {
			loc, err = time.LoadLocation(sched.Recurring.Timezone)
		}
		if err != nil || loc == nil {
			loc = time.Local
		}
		
		localNow := now.In(loc)
		
		// Check day of week
		if len(sched.Recurring.Days) > 0 {
			dayMatch := false
			for _, d := range sched.Recurring.Days {
				if localNow.Weekday() == d {
					dayMatch = true
					break
				}
			}
			if !dayMatch {
				return false
			}
		}

		// Check time of day (HH:MM format)
		currentTime := localNow.Format("15:04")
		if sched.Recurring.StartTime != "" && currentTime < sched.Recurring.StartTime {
			return false
		}
		if sched.Recurring.EndTime != "" && currentTime >= sched.Recurring.EndTime {
			return false
		}
	}

	return true
}

// Scheduler background service
type Scheduler struct {
	engine   *Engine
	interval time.Duration
	stopC    chan struct{}
}

func NewScheduler(eng *Engine, interval time.Duration) *Scheduler {
	return &Scheduler{
		engine:   eng,
		interval: interval,
		stopC:    make(chan struct{}),
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.evaluate()
		case <-ctx.Done():
			return
		case <-s.stopC:
			return
		}
	}
}

func (s *Scheduler) evaluate() {
	triggerReapply := false
	rs := s.engine.CurrentRules()
	if rs == nil {
		return
	}

	// 1. Check time-based schedules
	now := time.Now()
	for _, rule := range rs.Rules {
		if rule.Schedule != nil {
			if IsActive(rule.Schedule, now) {
				// Simplified check for demo
				triggerReapply = true
				break
			}
		}
	}

	// 2. Check dynamic IP sets (e.g. K8s pod changes)
	// In a real implementation, we'd cache the previous results of ResolveDynamicVars
	// and only trigger reapply if the IP list has changed.
	// For production readiness, per-second polling is heavy, so we would use 
	// K8s Watch API or a longer interval.
	
	if triggerReapply {
		log.Printf("Scheduler: triggering dynamic rule reapply")
		if err := s.engine.ReapplyRules(context.Background()); err != nil {
			log.Printf("Scheduler: failed to reapply rules: %v", err)
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopC)
}
