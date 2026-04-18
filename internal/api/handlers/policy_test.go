package handlers

import (
	"bytes"
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

func TestPolicyHandlers(t *testing.T) {
	// Setup
	tmpDir, _ := os.MkdirTemp("", "rampart-policy-api-test-*")
	defer os.RemoveAll(tmpDir)

	mockBE := &backend.MockBackend{}
	eng := engine.NewEngine(mockBE)
	snapStore, _ := snapshot.NewStore(tmpDir)
	auditStore, _ := audit.NewStore(tmpDir, time.Hour)
	handler := NewPolicyHandler(eng, snapStore, auditStore)

	t.Run("HandleCurrent_Empty", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/policies/current", nil)
		rr := httptest.NewRecorder()
		handler.HandleCurrent(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["status"] != "success" {
			t.Errorf("handler returned wrong status: got %v want %v", resp["status"], "success")
		}
	})

	t.Run("HandleApply", func(t *testing.T) {
		psBody := model.PolicySetYAML{
			APIVersion: "rampartfw.com/v1",
			Kind:       "PolicySet",
			Metadata:   model.PolicyMetadata{Name: "test-policy"},
			Policies: []model.PolicyYAML{
				{
					Name:     "p1",
					Priority: 100,
					Rules: []model.RuleYAML{
						{
							Name:   "r1",
							Action: model.ActionAccept,
							Match: model.MatchYAML{
								Protocol:  "tcp",
								DestPorts: 80,
							},
						},
					},
				},
			},
		}
		body, _ := json.Marshal(psBody)
		req, _ := http.NewRequest("POST", "/api/v1/policies/apply", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler.HandleApply(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if len(eng.CurrentRules().Rules) != 1 {
			t.Errorf("expected 1 rule in engine, got %d", len(eng.CurrentRules().Rules))
		}
	})

	t.Run("HandleSimulate", func(t *testing.T) {
		pkt := model.SimulatedPacket{
			SourceIP:  []byte{10, 0, 0, 1},
			DestIP:    []byte{192, 168, 1, 1},
			Protocol:  model.ProtocolTCP,
			DestPort:  80,
			Direction: model.DirectionInbound,
		}
		body, _ := json.Marshal(pkt)
		req, _ := http.NewRequest("POST", "/api/v1/policies/simulate", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler.HandleSimulate(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var resp struct {
			Status string                 `json:"status"`
			Data   model.SimulationResult `json:"data"`
		}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Data.Verdict != model.ActionAccept {
			t.Errorf("expected verdict ACCEPT, got %v", resp.Data.Verdict)
		}
	})

	t.Run("HandleFlush", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/policies", nil)
		rr := httptest.NewRecorder()
		handler.HandleFlush(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if len(eng.CurrentRules().Rules) != 0 {
			t.Errorf("expected 0 rules in engine after flush, got %d", len(eng.CurrentRules().Rules))
		}
	})
}

