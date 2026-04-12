package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rampartfw/rampart/internal/audit"
	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/model"
)

func TestRulesHandlers(t *testing.T) {
	// Setup
	tmpDir, _ := os.MkdirTemp("", "rampart-rules-test-*")
	defer os.RemoveAll(tmpDir)

	mockBE := &backend.MockBackend{}
	eng := engine.NewEngine(mockBE)
	auditStore, _ := audit.NewStore(tmpDir, time.Hour)
	handler := NewRulesHandler(eng, auditStore)

	t.Run("HandleAdd", func(t *testing.T) {
		rule := model.CompiledRule{
			Name:     "test-rule",
			Priority: 100,
			Action:   model.ActionAccept,
		}
		body, _ := json.Marshal(rule)
		req, _ := http.NewRequest("POST", "/api/v1/rules", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler.HandleAdd(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		if len(eng.CurrentRules().Rules) != 1 {
			t.Errorf("expected 1 rule in engine, got %d", len(eng.CurrentRules().Rules))
		}
	})

	t.Run("HandleList", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/rules", nil)
		rr := httptest.NewRecorder()
		handler.HandleList(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var resp struct {
			Status string               `json:"status"`
			Data   []model.CompiledRule `json:"data"`
		}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if len(resp.Data) != 1 {
			t.Errorf("expected 1 rule in list, got %d", len(resp.Data))
		}
	})

	t.Run("HandleGet", func(t *testing.T) {
		ruleID := eng.CurrentRules().Rules[0].ID
		req, _ := http.NewRequest("GET", "/api/v1/rules/"+ruleID, nil)
		
		// Add ID to context since we are bypassing the router
		ctx := context.WithValue(req.Context(), ParamsKey, map[string]string{"id": ruleID})
		req = req.WithContext(ctx)
		
		rr := httptest.NewRecorder()
		handler.HandleGet(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("HandleDelete", func(t *testing.T) {
		ruleID := eng.CurrentRules().Rules[0].ID
		req, _ := http.NewRequest("DELETE", "/api/v1/rules/"+ruleID, nil)
		
		ctx := context.WithValue(req.Context(), ParamsKey, map[string]string{"id": ruleID})
		req = req.WithContext(ctx)
		
		rr := httptest.NewRecorder()
		handler.HandleDelete(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if len(eng.CurrentRules().Rules) != 0 {
			t.Errorf("expected 0 rules in engine after delete, got %d", len(eng.CurrentRules().Rules))
		}
	})
}
