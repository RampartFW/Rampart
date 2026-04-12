package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/rampartfw/rampart/internal/audit"
	"github.com/rampartfw/rampart/internal/model"
)

type AuditHandler struct {
	store *audit.Store
}

func NewAuditHandler(as *audit.Store) *AuditHandler {
	return &AuditHandler{
		store: as,
	}
}

func (h *AuditHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = val
		}
	}

	query := audit.AuditQuery{
		Limit:  limit,
		Offset: offset,
	}

	events, total, err := h.store.Search(query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"events": events,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *AuditHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	params := Params(r)
	id := params["id"]
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing audit ID")
		return
	}

	event, err := h.store.Get(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Audit event not found")
		return
	}

	respondJSON(w, http.StatusOK, event)
}

func (h *AuditHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	
	query := audit.AuditQuery{
		Action:   model.AuditAction(q.Get("action")),
		Actor:    q.Get("actor"),
		Resource: q.Get("resource"),
	}

	if s := q.Get("since"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			query.Since = t
		}
	}

	if u := q.Get("until"); u != "" {
		if t, err := time.Parse(time.RFC3339, u); err == nil {
			query.Until = t
		}
	}

	if l := q.Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			query.Limit = val
		}
	} else {
		query.Limit = 50
	}

	if o := q.Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			query.Offset = val
		}
	}

	events, total, err := h.store.Search(query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"events": events,
		"total":  total,
		"limit":  query.Limit,
		"offset": query.Offset,
	})
}
