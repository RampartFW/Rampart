package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/model"
	"github.com/rampartfw/rampart/internal/snapshot"
	"github.com/rampartfw/rampart/internal/audit"
)

type SnapshotHandler struct {
	store      *snapshot.Store
	engine     *engine.Engine
	auditStore *audit.Store
}

func NewSnapshotHandler(s *snapshot.Store, eng *engine.Engine, as *audit.Store) *SnapshotHandler {
	return &SnapshotHandler{
		store:      s,
		engine:     eng,
		auditStore: as,
	}
}

func (h *SnapshotHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	snaps, err := h.store.List()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, snaps)
}

func (h *SnapshotHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Trigger     string `json:"trigger"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	current := h.engine.CurrentRules()
	if current == nil {
		current = &model.CompiledRuleSet{}
	}

	snap, err := h.store.Create(req.Trigger, req.Description, current)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
// Record audit event
h.auditStore.Record(model.AuditEvent{
	Action: model.AuditSnapshot,
	Actor:  model.AuditActor{Type: "api", Identity: r.RemoteAddr},
	Result: model.AuditResultSuccess,
})


	respondJSON(w, http.StatusCreated, snap)
}

func (h *SnapshotHandler) HandleRollback(w http.ResponseWriter, r *http.Request) {
	params := Params(r)
	id := params["id"]
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing snapshot ID")
		return
	}

	_, snapState, err := h.store.Load(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Snapshot not found")
		return
	}

	if err := h.engine.Backend().Apply(snapState); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.engine.SetRules(snapState)

	h.auditStore.Record(model.AuditEvent{
		Action: model.AuditRollback,
		Actor:  model.AuditActor{Type: "api", Identity: r.RemoteAddr},
		Result: model.AuditResultSuccess,
	})

	respondJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (h *SnapshotHandler) HandleDiff(w http.ResponseWriter, r *http.Request) {
	params := Params(r)
	id := params["id"]
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing snapshot ID")
		return
	}

	current := h.engine.CurrentRules()
	if current == nil {
		current = &model.CompiledRuleSet{}
	}

	plan, err := h.store.Diff(id, current)
	if err != nil {
		respondError(w, http.StatusNotFound, "Snapshot not found")
		return
	}

	respondJSON(w, http.StatusOK, plan)
}

func (h *SnapshotHandler) HandleExport(w http.ResponseWriter, r *http.Request) {
	params := Params(r)
	id := params["id"]
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing snapshot ID")
		return
	}

	_, snapState, err := h.store.Load(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Snapshot not found")
		return
	}

	// For export, we return the CompiledRuleSet which can be converted to YAML on client side
	// or we could provide a helper to convert to PolicySetYAML.
	// For now, let's just return the CompiledRuleSet.
	respondJSON(w, http.StatusOK, snapState)
}

func (h *SnapshotHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	params := Params(r)
	id := params["id"]
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing snapshot ID")
		return
	}

	if err := h.store.Delete(id); err != nil {
		respondError(w, http.StatusNotFound, "Snapshot not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"id": id})
}
