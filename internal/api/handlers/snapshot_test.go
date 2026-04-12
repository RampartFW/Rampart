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
	"github.com/rampartfw/rampart/internal/snapshot"
)

func TestSnapshotHandlers(t *testing.T) {
	// Setup
	tmpDir, _ := os.MkdirTemp("", "rampart-snap-test-*")
	defer os.RemoveAll(tmpDir)

	mockBE := &backend.MockBackend{}
	eng := engine.NewEngine(mockBE)
	snapStore, _ := snapshot.NewStore(tmpDir)
	auditStore, _ := audit.NewStore(tmpDir, time.Hour)
	handler := NewSnapshotHandler(snapStore, eng, auditStore)

	var snapID string

	t.Run("HandleCreate", func(t *testing.T) {
		reqBody := map[string]string{
			"trigger":     "manual",
			"description": "test snapshot",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/v1/snapshots", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler.HandleCreate(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var resp struct {
			Status string         `json:"status"`
			Data   model.Snapshot `json:"data"`
		}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		snapID = resp.Data.ID
		if snapID == "" {
			t.Fatal("expected snapshot ID in response")
		}
	})

	t.Run("HandleList", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/snapshots", nil)
		rr := httptest.NewRecorder()
		handler.HandleList(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var resp struct {
			Status string           `json:"status"`
			Data   []model.Snapshot `json:"data"`
		}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if len(resp.Data) != 1 {
			t.Errorf("expected 1 snapshot in list, got %d", len(resp.Data))
		}
	})

	t.Run("HandleDiff", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/snapshots/"+snapID+"/diff", nil)
		ctx := context.WithValue(req.Context(), ParamsKey, map[string]string{"id": snapID})
		req = req.WithContext(ctx)
		
		rr := httptest.NewRecorder()
		handler.HandleDiff(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("HandleRollback", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/snapshots/"+snapID+"/rollback", nil)
		ctx := context.WithValue(req.Context(), ParamsKey, map[string]string{"id": snapID})
		req = req.WithContext(ctx)
		
		rr := httptest.NewRecorder()
		handler.HandleRollback(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("HandleDelete", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/snapshots/"+snapID, nil)
		ctx := context.WithValue(req.Context(), ParamsKey, map[string]string{"id": snapID})
		req = req.WithContext(ctx)
		
		rr := httptest.NewRecorder()
		handler.HandleDelete(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		snaps, _ := snapStore.List()
		if len(snaps) != 0 {
			t.Errorf("expected 0 snapshots after delete, got %d", len(snaps))
		}
	})
}
