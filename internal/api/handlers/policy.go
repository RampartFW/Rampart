package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/snapshot"
	"github.com/rampartfw/rampart/internal/audit"
	"github.com/rampartfw/rampart/internal/model"
)

type PolicyHandler struct {
	engine        *engine.Engine
	snapshotStore *snapshot.Store
	auditStore    *audit.Store
}

func NewPolicyHandler(eng *engine.Engine, ss *snapshot.Store, as *audit.Store) *PolicyHandler {
	return &PolicyHandler{
		engine:        eng,
		snapshotStore: ss,
		auditStore:    as,
	}
}

func (h *PolicyHandler) HandlePlan(w http.ResponseWriter, r *http.Request) {
	var ps model.PolicySetYAML
	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	compiled, err := engine.Compile(&ps, nil)
	if err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	current := h.engine.CurrentRules()
	if current == nil {
		current = &model.CompiledRuleSet{}
	}

	plan := engine.GeneratePlan(current, compiled)
	respondJSON(w, http.StatusOK, plan)
}

func (h *PolicyHandler) HandleApply(w http.ResponseWriter, r *http.Request) {
	var ps model.PolicySetYAML
	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	compiled, err := engine.Compile(&ps, nil)
	if err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	current := h.engine.CurrentRules()
	var beforeJSON []byte
	if current != nil {
		beforeJSON, _ = json.Marshal(current)
		h.snapshotStore.Create("pre-apply", "Auto snapshot before policy apply", current)
	}

	if err := h.engine.Backend().Apply(r.Context(), compiled); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.engine.SetRules(compiled)
	afterJSON, _ := json.Marshal(compiled)

	// Record audit event with full state
	h.auditStore.Record(model.AuditEvent{
		Action: model.AuditApply,
		Actor:  model.AuditActor{Type: "api", Identity: r.RemoteAddr},
		Resource: model.AuditResource{Type: "policy", Name: ps.Metadata.Name},
		Before: beforeJSON,
		After:  afterJSON,
		Result: model.AuditResultSuccess,
	})

	h.engine.Broadcast(model.AuditEvent{
		Action: model.AuditApply,
		Actor:  model.AuditActor{Type: "api", Identity: r.RemoteAddr},
		Result: model.AuditResultSuccess,
	})

	respondJSON(w, http.StatusOK, map[string]string{"message": "Policy applied successfully"})
}

func (h *PolicyHandler) HandleSimulate(w http.ResponseWriter, r *http.Request) {
	var pkt model.SimulatedPacket
	if err := json.NewDecoder(r.Body).Decode(&pkt); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	current := h.engine.CurrentRules()
	if current == nil {
		respondError(w, http.StatusNotFound, "No rules active")
		return
	}

	result := engine.Simulate(current.Rules, pkt)
	respondJSON(w, http.StatusOK, result)
}

func (h *PolicyHandler) HandleCurrent(w http.ResponseWriter, r *http.Request) {
	current := h.engine.CurrentRules()
	if current == nil {
		respondJSON(w, http.StatusOK, nil)
		return
	}
	respondJSON(w, http.StatusOK, current)
}

func (h *PolicyHandler) HandleConflicts(w http.ResponseWriter, r *http.Request) {
	current := h.engine.CurrentRules()
	if current == nil {
		respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	conflicts := engine.DetectConflicts(current.Rules)
	respondJSON(w, http.StatusOK, conflicts)
}

func (h *PolicyHandler) HandleFlush(w http.ResponseWriter, r *http.Request) {
	current := h.engine.CurrentRules()
	var beforeJSON []byte
	if current != nil {
		beforeJSON, _ = json.Marshal(current)
	}

	if err := h.engine.Backend().Flush(r.Context()); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.engine.SetRules(&model.CompiledRuleSet{})

	h.auditStore.Record(model.AuditEvent{
		Action: model.AuditFlush,
		Actor:  model.AuditActor{Type: "api", Identity: r.RemoteAddr},
		Before: beforeJSON,
		Result: model.AuditResultSuccess,
	})

	h.engine.Broadcast(model.AuditEvent{
		Action: model.AuditFlush,
		Actor:  model.AuditActor{Type: "api", Identity: r.RemoteAddr},
		Result: model.AuditResultSuccess,
	})

	respondJSON(w, http.StatusOK, map[string]string{"message": "All rules flushed"})
}
