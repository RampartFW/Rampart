package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/audit"
	"github.com/rampartfw/rampart/internal/model"
)

type RulesHandler struct {
	engine     *engine.Engine
	auditStore *audit.Store
}

func NewRulesHandler(eng *engine.Engine, as *audit.Store) *RulesHandler {
	return &RulesHandler{
		engine:     eng,
		auditStore: as,
	}
}

func (h *RulesHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	current := h.engine.CurrentRules()
	if current == nil {
		respondJSON(w, http.StatusOK, []interface{}{})
		return
	}
	respondJSON(w, http.StatusOK, current.Rules)
}

func (h *RulesHandler) HandleAdd(w http.ResponseWriter, r *http.Request) {
	var rule model.CompiledRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if rule.ID == "" {
		rule.ID = model.GenerateUUIDv7()
	}

	current := h.engine.CurrentRules()
	var beforeJSON []byte
	var newRules []model.CompiledRule
	if current != nil {
		beforeJSON, _ = json.Marshal(current)
		newRules = append(newRules, current.Rules...)
	}
	newRules = append(newRules, rule)

	newSet := &model.CompiledRuleSet{
		Rules: newRules,
	}

	if err := h.engine.Backend().Apply(r.Context(), newSet); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.engine.SetRules(newSet)
	afterJSON, _ := json.Marshal(newSet)

	h.auditStore.Record(model.AuditEvent{
		Action: model.AuditApply,
		Actor:  model.AuditActor{Type: "api", Identity: r.RemoteAddr},
		Resource: model.AuditResource{Type: "rule", ID: rule.ID, Name: rule.Name},
		Before: beforeJSON,
		After:  afterJSON,
		Result: model.AuditResultSuccess,
	})

	respondJSON(w, http.StatusCreated, rule)
}

func (h *RulesHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	params := Params(r)
	id := params["id"]
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing rule ID")
		return
	}

	current := h.engine.CurrentRules()
	if current == nil {
		respondError(w, http.StatusNotFound, "Rule not found")
		return
	}

	for _, rule := range current.Rules {
		if rule.ID == id {
			respondJSON(w, http.StatusOK, rule)
			return
		}
	}

	respondError(w, http.StatusNotFound, "Rule not found")
}

func (h *RulesHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	params := Params(r)
	id := params["id"]
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing rule ID")
		return
	}

	current := h.engine.CurrentRules()
	if current == nil {
		respondError(w, http.StatusNotFound, "Rule not found")
		return
	}

	found := false
	var ruleName string
	var newRules []model.CompiledRule
	for _, rule := range current.Rules {
		if rule.ID == id {
			found = true
			ruleName = rule.Name
			continue
		}
		newRules = append(newRules, rule)
	}

	if !found {
		respondError(w, http.StatusNotFound, "Rule not found")
		return
	}

	beforeJSON, _ := json.Marshal(current)
	newSet := &model.CompiledRuleSet{
		Rules: newRules,
	}

	if err := h.engine.Backend().Apply(r.Context(), newSet); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.engine.SetRules(newSet)
	afterJSON, _ := json.Marshal(newSet)

	h.auditStore.Record(model.AuditEvent{
		Action: model.AuditApply,
		Actor:  model.AuditActor{Type: "api", Identity: r.RemoteAddr},
		Resource: model.AuditResource{Type: "rule", ID: id, Name: ruleName},
		Before: beforeJSON,
		After:  afterJSON,
		Result: model.AuditResultSuccess,
	})

	respondJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (h *RulesHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.engine.Backend().Stats(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, stats)
}
