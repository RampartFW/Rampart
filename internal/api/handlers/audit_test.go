package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rampartfw/rampart/internal/audit"
	"github.com/rampartfw/rampart/internal/model"
)

func TestAuditHandlers(t *testing.T) {
	// Setup
	tmpDir, _ := os.MkdirTemp("", "rampart-audit-api-test-*")
	defer os.RemoveAll(tmpDir)

	auditStore, _ := audit.NewStore(tmpDir, time.Hour)
	handler := NewAuditHandler(auditStore)

	// Inject an event
	event := model.AuditEvent{
		Action: model.AuditApply,
		Actor:  model.AuditActor{Type: "user", Identity: "test-user"},
	}
	auditStore.Record(event)
	
	// Wait for writer goroutine to process the event
	time.Sleep(100 * time.Millisecond)

	t.Run("HandleList", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/audit", nil)
		rr := httptest.NewRecorder()
		handler.HandleList(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var resp struct {
			Status string `json:"status"`
			Data   struct {
				Events []model.AuditEvent `json:"events"`
				Total  int                `json:"total"`
			} `json:"data"`
		}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Data.Total != 1 {
			t.Errorf("expected 1 audit event, got %d", resp.Data.Total)
		}
	})

	t.Run("HandleGet", func(t *testing.T) {
		query := audit.AuditQuery{Limit: 1}
		events, _, _ := auditStore.Search(query)
		eventID := events[0].ID

		req, _ := http.NewRequest("GET", "/api/v1/audit/"+eventID, nil)
		ctx := context.WithValue(req.Context(), ParamsKey, map[string]string{"id": eventID})
		req = req.WithContext(ctx)
		
		rr := httptest.NewRecorder()
		handler.HandleGet(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("HandleSearch", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/audit/search?actor=test-user", nil)
		rr := httptest.NewRecorder()
		handler.HandleSearch(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var resp struct {
			Status string `json:"status"`
			Data   struct {
				Events []model.AuditEvent `json:"events"`
				Total  int                `json:"total"`
			} `json:"data"`
		}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Data.Total != 1 {
			t.Errorf("expected 1 audit event in search, got %d", resp.Data.Total)
		}
	})
}
