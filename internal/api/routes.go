package api

import (
	"net/http"
	"github.com/rampartfw/rampart/internal/api/handlers"
)

func (s *Server) routes() {
	// Standard middleware (apply to all)
	s.router.Use(RecoveryMiddleware)
	s.router.Use(RequestIDMiddleware)
	s.router.Use(LoggingMiddleware)
	s.router.Use(CORSMiddleware)

	// Create an auth-wrapped handler for API routes
	apiAuth := AuthMiddleware(s.cfg.API.Keys)

	// Policy handlers
	ph := handlers.NewPolicyHandler(s.engine, s.snapshotStore, s.auditStore, s.raftNode)
	s.router.Handle("POST", "/api/v1/policies/plan", apiAuth(http.HandlerFunc(ph.HandlePlan)).ServeHTTP)
	s.router.Handle("POST", "/api/v1/policies/apply", apiAuth(http.HandlerFunc(ph.HandleApply)).ServeHTTP)
	s.router.Handle("POST", "/api/v1/policies/simulate", apiAuth(http.HandlerFunc(ph.HandleSimulate)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/policies/current", apiAuth(http.HandlerFunc(ph.HandleCurrent)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/policies/conflicts", apiAuth(http.HandlerFunc(ph.HandleConflicts)).ServeHTTP)
	s.router.Handle("DELETE", "/api/v1/policies", apiAuth(http.HandlerFunc(ph.HandleFlush)).ServeHTTP)

	// Rules handlers
	// Note: RulesHandler might also need RaftNode if we want quick-add rules to be replicated.
	// For now, only PolicySet updates are replicated.
	rh := handlers.NewRulesHandler(s.engine, s.auditStore)
	s.router.Handle("GET", "/api/v1/rules", apiAuth(http.HandlerFunc(rh.HandleList)).ServeHTTP)
	s.router.Handle("POST", "/api/v1/rules", apiAuth(http.HandlerFunc(rh.HandleAdd)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/rules/:id", apiAuth(http.HandlerFunc(rh.HandleGet)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/rules/:id/impact", apiAuth(http.HandlerFunc(rh.HandleImpact)).ServeHTTP)
	s.router.Handle("DELETE", "/api/v1/rules/:id", apiAuth(http.HandlerFunc(rh.HandleDelete)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/rules/:id/stats", apiAuth(http.HandlerFunc(rh.HandleStats)).ServeHTTP)

	// Snapshot handlers
	sh := handlers.NewSnapshotHandler(s.snapshotStore, s.engine, s.auditStore)
	s.router.Handle("GET", "/api/v1/snapshots", apiAuth(http.HandlerFunc(sh.HandleList)).ServeHTTP)
	s.router.Handle("POST", "/api/v1/snapshots", apiAuth(http.HandlerFunc(sh.HandleCreate)).ServeHTTP)
	s.router.Handle("POST", "/api/v1/snapshots/:id/rollback", apiAuth(http.HandlerFunc(sh.HandleRollback)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/snapshots/:id/diff", apiAuth(http.HandlerFunc(sh.HandleDiff)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/snapshots/:id/export", apiAuth(http.HandlerFunc(sh.HandleExport)).ServeHTTP)
	s.router.Handle("DELETE", "/api/v1/snapshots/:id", apiAuth(http.HandlerFunc(sh.HandleDelete)).ServeHTTP)

	// Audit handlers
	ah := handlers.NewAuditHandler(s.auditStore)
	s.router.Handle("GET", "/api/v1/audit", apiAuth(http.HandlerFunc(ah.HandleList)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/audit/verify", apiAuth(http.HandlerFunc(ah.HandleVerify)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/audit/:id", apiAuth(http.HandlerFunc(ah.HandleGet)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/audit/search", apiAuth(http.HandlerFunc(ah.HandleSearch)).ServeHTTP)

	// System handlers
	sh_sys := handlers.NewSystemHandler(s.cfg, s.engine, s.raftNode)
	s.router.Handle("GET", "/api/v1/system/info", apiAuth(http.HandlerFunc(sh_sys.HandleInfo)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/system/health", apiAuth(http.HandlerFunc(sh_sys.HandleHealth)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/system/backends", apiAuth(http.HandlerFunc(sh_sys.HandleBackends)).ServeHTTP)
	s.router.Handle("GET", "/api/v1/cluster/status", apiAuth(http.HandlerFunc(sh_sys.HandleClusterStatus)).ServeHTTP)
	s.router.Handle("GET", "/metrics", sh_sys.HandleMetrics) // Usually internal, can be public or separate auth

	// SSE
	s.router.Handle("GET", "/api/v1/events", apiAuth(http.HandlerFunc(s.HandleSSE)).ServeHTTP)

	// WebUI (No Auth)
	s.router.Handle("GET", "/ui/", s.serveUI().ServeHTTP)
	s.router.Handle("GET", "/ui/*", s.serveUI().ServeHTTP)
}
