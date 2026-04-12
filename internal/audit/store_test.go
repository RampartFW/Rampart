package audit

import (
	"os"
	"testing"
	"time"

	"github.com/rampartfw/rampart/internal/model"
)

func TestAuditStore(t *testing.T) {
	dir, err := os.MkdirTemp("", "rampart-audit-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	store, err := NewStore(dir, 90*24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	event := model.AuditEvent{
		Action: model.AuditApply,
		Actor: model.AuditActor{
			Type:     "user",
			Identity: "tester",
		},
		Timestamp: time.Now(),
		Result: model.AuditResult{
			Status: "success",
		},
	}

	// Record
	err = store.Record(event)
	if err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	// Give the writer some time
	time.Sleep(100 * time.Millisecond)

	// Search
	results, total, err := store.Search(AuditQuery{Actor: "tester"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if total != 1 {
		t.Errorf("Expected 1 result, got %d", total)
	}
	if results[0].Actor.Identity != "tester" {
		t.Errorf("Expected identity tester, got %s", results[0].Actor.Identity)
	}

	// Get
	event2, err := store.Get(results[0].ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if event2.ID != results[0].ID {
		t.Errorf("Expected ID %s, got %s", results[0].ID, event2.ID)
	}

	// Integrity
	ok, err := store.VerifyIntegrity()
	if err != nil {
		t.Fatalf("VerifyIntegrity failed: %v", err)
	}
	if !ok {
		t.Error("Integrity check failed")
	}

	// Test multi-event chain
	for i := 0; i < 5; i++ {
		_ = store.Record(model.AuditEvent{
			Action: model.AuditSnapshot,
			Actor: model.AuditActor{Identity: "chain-tester"},
		})
	}
	time.Sleep(100 * time.Millisecond)

	ok, err = store.VerifyIntegrity()
	if err != nil {
		t.Fatalf("VerifyIntegrity failed after multi-event: %v", err)
	}
	if !ok {
		t.Error("Integrity check failed after multi-event")
	}
}
