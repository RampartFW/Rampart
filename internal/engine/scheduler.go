package engine

import (
	"context"
	"time"

	"github.com/rampartfw/rampart/internal/model"
)

// IsActive checks if a schedule is active at the given time.
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
		loc := time.Local
		if sched.Recurring.Timezone != "" {
			var err error
			loc, err = time.LoadLocation(sched.Recurring.Timezone)
			if err != nil {
				// Fallback to local if timezone is invalid
				loc = time.Local
			}
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

		// Check time of day
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

type Scheduler struct {
	engine   *Engine
	interval time.Duration
}

func NewScheduler(engine *Engine, interval time.Duration) *Scheduler {
	if interval == 0 {
		interval = 30 * time.Second
	}
	return &Scheduler{
		engine:   engine,
		interval: interval,
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
		}
	}
}

func (s *Scheduler) evaluate() {
	now := time.Now()
	rules := s.engine.CurrentRules()
	if rules == nil {
		return
	}
	changed := false

	for i := range rules.Rules {
		rule := &rules.Rules[i]
		if rule.Schedule == nil {
			continue
		}

		isActive := IsActive(rule.Schedule, now)
		if isActive != rule.Schedule.WasActive {
			rule.Schedule.WasActive = isActive
			changed = true
		}
	}

	if changed {
		s.engine.ReapplyRules(context.Background())
	}
}
