package api

import (
	"github.com/rampartfw/rampart/internal/api/handlers"
)

func (s *Server) routes() {
	// Standard middleware
	s.router.Use(RequestIDMiddleware)
	s.router.Use(LoggingMiddleware)
	s.router.Use(CORSMiddleware)
	s.router.Use(AuthMiddleware(s.cfg.API.Keys))

	// Policy handlers
	ph := handlers.NewPolicyHandler(s.engine, s.snapshotStore, s.auditStore)
	s.router.Handle("POST", "/api/v1/policies/plan", ph.HandlePlan)
	s.router.Handle("POST", "/api/v1/policies/apply", ph.HandleApply)
	s.router.Handle("POST", "/api/v1/policies/simulate", ph.HandleSimulate)
	s.router.Handle("GET", "/api/v1/policies/current", ph.HandleCurrent)
	s.router.Handle("GET", "/api/v1/policies/conflicts", ph.HandleConflicts)
	s.router.Handle("DELETE", "/api/v1/policies", ph.HandleFlush)

	// Rules handlers
	rh := handlers.NewRulesHandler(s.engine, s.auditStore)
	s.router.Handle("GET", "/api/v1/rules", rh.HandleList)
	s.router.Handle("POST", "/api/v1/rules", rh.HandleAdd)
	s.router.Handle("GET", "/api/v1/rules/:id", rh.HandleGet)
	s.router.Handle("DELETE", "/api/v1/rules/:id", rh.HandleDelete)
	s.router.Handle("GET", "/api/v1/rules/:id/stats", rh.HandleStats)

	// Snapshot handlers
	sh := handlers.NewSnapshotHandler(s.snapshotStore, s.engine, s.auditStore)
	s.router.Handle("GET", "/api/v1/snapshots", sh.HandleList)
	s.router.Handle("POST", "/api/v1/snapshots", sh.HandleCreate)
	s.router.Handle("POST", "/api/v1/snapshots/:id/rollback", sh.HandleRollback)
	s.router.Handle("GET", "/api/v1/snapshots/:id/diff", sh.HandleDiff)
	s.router.Handle("GET", "/api/v1/snapshots/:id/export", sh.HandleExport)
	s.router.Handle("DELETE", "/api/v1/snapshots/:id", sh.HandleDelete)

	// Audit handlers
	ah := handlers.NewAuditHandler(s.auditStore)
	s.router.Handle("GET", "/api/v1/audit", ah.HandleList)
	s.router.Handle("GET", "/api/v1/audit/:id", ah.HandleGet)
	s.router.Handle("GET", "/api/v1/audit/search", ah.HandleSearch)

	// System handlers
	sh_sys := handlers.NewSystemHandler(s.cfg, s.engine)
	s.router.Handle("GET", "/api/v1/system/info", sh_sys.HandleInfo)
	s.router.Handle("GET", "/api/v1/system/health", sh_sys.HandleHealth)
	s.router.Handle("GET", "/api/v1/system/backends", sh_sys.HandleBackends)
	s.router.Handle("GET", "/metrics", sh_sys.HandleMetrics)

	// SSE
	s.router.Handle("GET", "/api/v1/events", s.HandleSSE)

	// WebUI
	s.router.Handle("GET", "/ui/*", s.serveUI().ServeHTTP)
}
