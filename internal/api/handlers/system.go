package handlers

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/rampartfw/rampart/internal/config"
	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/model"
)

type SystemHandler struct {
	cfg       *config.Config
	engine    *engine.Engine
	raftNode  RaftNode
	startTime time.Time
}

func NewSystemHandler(cfg *config.Config, engine *engine.Engine, raftNode RaftNode) *SystemHandler {
	return &SystemHandler{
		cfg:       cfg,
		engine:    engine,
		raftNode:  raftNode,
		startTime: time.Now(),
	}
}

func (h *SystemHandler) HandleClusterStatus(w http.ResponseWriter, r *http.Request) {
	if h.raftNode == nil {
		respondError(w, http.StatusServiceUnavailable, "Clustering is not enabled on this node")
		return
	}

	status := h.raftNode.Status()
	respondJSON(w, http.StatusOK, status)
}

func (h *SystemHandler) HandleInfo(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"version": "1.0.0-draft",
		"go_version": runtime.Version(),
		"arch":      runtime.GOARCH,
		"os":        runtime.GOOS,
		"uptime":    time.Since(h.startTime).String(),
		"config": map[string]interface{}{
			"loaded_from": h.cfg.LoadedFrom(),
			"backend":     h.cfg.Backend.Type,
			"cluster":     h.cfg.Cluster.Enabled,
		},
	}
	respondJSON(w, http.StatusOK, info)
}

func (h *SystemHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status": "healthy",
		"checks": map[string]string{
			"api":     "ok",
			"backend": "ok",
		},
	}
	if h.engine.Backend() == nil {
		health["status"] = "unhealthy"
		health["checks"].(map[string]string)["backend"] = "error: no backend"
	} else if err := h.engine.Backend().Probe(); err != nil {
		health["status"] = "degraded"
		health["checks"].(map[string]string)["backend"] = fmt.Sprintf("error: %v", err)
	}

	respondJSON(w, http.StatusOK, health)
}

func (h *SystemHandler) HandleBackends(w http.ResponseWriter, r *http.Request) {
	backends := []map[string]interface{}{
		{
			"name":    "nftables",
			"active":  h.cfg.Backend.Type == "nftables" || (h.cfg.Backend.Type == "auto" && h.engine.Backend().Name() == "nftables"),
			"support": "linux",
		},
		{
			"name":    "iptables",
			"active":  h.cfg.Backend.Type == "iptables" || (h.cfg.Backend.Type == "auto" && h.engine.Backend().Name() == "iptables"),
			"support": "linux",
		},
		{
			"name":    "ebpf",
			"active":  h.cfg.Backend.Type == "ebpf" || (h.cfg.Backend.Type == "auto" && h.engine.Backend().Name() == "ebpf"),
			"support": "linux (XDP)",
		},
	}
	respondJSON(w, http.StatusOK, backends)
}

func (h *SystemHandler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	// Simple Prometheus text format
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	
	fmt.Fprintf(w, "# HELP rampart_uptime_seconds Uptime in seconds\n")
	fmt.Fprintf(w, "# TYPE rampart_uptime_seconds gauge\n")
	fmt.Fprintf(w, "rampart_uptime_seconds %f\n", time.Since(h.startTime).Seconds())

	fmt.Fprintf(w, "# HELP rampart_goroutines Number of goroutines\n")
	fmt.Fprintf(w, "# TYPE rampart_goroutines gauge\n")
	fmt.Fprintf(w, "rampart_goroutines %d\n", runtime.NumGoroutine())

	if h.engine.CurrentRules() != nil {
		fmt.Fprintf(w, "# HELP rampart_rules_count Number of active rules\n")
		fmt.Fprintf(w, "# TYPE rampart_rules_count gauge\n")
		fmt.Fprintf(w, "rampart_rules_count %d\n", len(h.engine.CurrentRules().Rules))
	}

	// Cluster Metrics
	if h.raftNode != nil {
		status := h.raftNode.Status()
		fmt.Fprintf(w, "# HELP rampart_cluster_state Raft node state (0=Follower, 1=Candidate, 2=Leader)\n")
		fmt.Fprintf(w, "# TYPE rampart_cluster_state gauge\n")
		stateVal := 0
		switch status.State {
		case model.StateCandidate:
			stateVal = 1
		case model.StateLeader:
			stateVal = 2
		}
		fmt.Fprintf(w, "rampart_cluster_state %d\n", stateVal)
		
		fmt.Fprintf(w, "# HELP rampart_cluster_healthy Cluster health status\n")
		fmt.Fprintf(w, "# TYPE rampart_cluster_healthy gauge\n")
		healthyVal := 1
		if !status.IsHealthy {
			healthyVal = 0
		}
		fmt.Fprintf(w, "rampart_cluster_healthy %d\n", healthyVal)
	}

	// Backend Stats
	stats, err := h.engine.Backend().Stats(r.Context())
	if err == nil {
		fmt.Fprintf(w, "# HELP rampart_rule_packets_total Total packets matched by rule\n")
		fmt.Fprintf(w, "# TYPE rampart_rule_packets_total counter\n")
		fmt.Fprintf(w, "# HELP rampart_rule_bytes_total Total bytes matched by rule\n")
		fmt.Fprintf(w, "# TYPE rampart_rule_bytes_total counter\n")

		for id, s := range stats {
			// Find rule name for label if possible
			ruleName := id
			if h.engine.CurrentRules() != nil {
				for _, r := range h.engine.CurrentRules().Rules {
					if r.ID == id {
						ruleName = r.Name
						break
					}
				}
			}
			fmt.Fprintf(w, "rampart_rule_packets_total{rule_id=\"%s\", rule_name=\"%s\"} %d\n", id, ruleName, s.Packets)
			fmt.Fprintf(w, "rampart_rule_bytes_total{rule_id=\"%s\", rule_name=\"%s\"} %d\n", id, ruleName, s.Bytes)
		}
	}
}
