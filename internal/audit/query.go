package audit

import (
	"time"

	"github.com/rampartfw/rampart/internal/model"
)

// AuditQuery defines filters for searching audit events.
type AuditQuery struct {
	Action    model.AuditAction
	Actor     string
	Since     time.Time
	Until     time.Time
	Resource  string
	Limit     int
	Offset    int
}

// Matches returns true if the event matches the query filters.
func (q AuditQuery) Matches(event model.AuditEvent) bool {
	if q.Action != "" && event.Action != q.Action {
		return false
	}
	if q.Actor != "" && event.Actor.Identity != q.Actor {
		return false
	}
	if q.Resource != "" && event.Resource.Type != q.Resource {
		return false
	}
	if !q.Since.IsZero() && event.Timestamp.Before(q.Since) {
		return false
	}
	if !q.Until.IsZero() && event.Timestamp.After(q.Until) {
		return false
	}
	return true
}
