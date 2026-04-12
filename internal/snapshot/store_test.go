package snapshot

import (
	"os"
	"testing"
	"time"

	"github.com/rampartfw/rampart/internal/model"
)

func TestSnapshotStore(t *testing.T) {
	dir, err := os.MkdirTemp("", "rampart-snapshot-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	store, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	state := &model.CompiledRuleSet{
		Hash: "test-hash",
		Rules: []model.CompiledRule{
			{Name: "rule1", Priority: 10},
		},
		Backend: "nftables",
	}

	// Create
	snap, err := store.Create("manual", "test snapshot", state)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if snap.Description != "test snapshot" {
		t.Errorf("Expected description 'test snapshot', got '%s'", snap.Description)
	}

	// List
	snaps, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(snaps) != 1 {
		t.Errorf("Expected 1 snapshot, got %d", len(snaps))
	}

	// Load
	snap2, state2, err := store.Load(snap.ID)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if snap2.ID != snap.ID {
		t.Errorf("Expected ID %s, got %s", snap.ID, snap2.ID)
	}
	if state2.Hash != state.Hash {
		t.Errorf("Expected hash %s, got %s", state.Hash, state2.Hash)
	}
	if len(state2.Rules) != 1 || state2.Rules[0].Name != "rule1" {
		t.Errorf("Rules not loaded correctly")
	}

	// Diff
	current := &model.CompiledRuleSet{
		Hash: "new-hash",
		Rules: []model.CompiledRule{
			{Name: "rule1", Priority: 10},
			{Name: "rule2", Priority: 20},
		},
	}
	plan, err := store.Diff(snap.ID, current)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}
	if plan.AddCount != 1 {
		t.Errorf("Expected 1 add, got %d", plan.AddCount)
	}

	// Delete
	err = store.Delete(snap.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	snaps, _ = store.List()
	if len(snaps) != 0 {
		t.Errorf("Expected 0 snapshots after delete, got %d", len(snaps))
	}
}

func TestSnapshotRetention(t *testing.T) {
	dir, err := os.MkdirTemp("", "rampart-retention-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	store, _ := NewStore(dir)
	state := &model.CompiledRuleSet{Hash: "h"}

	// Create 5 snapshots
	for i := 0; i < 5; i++ {
		_, _ = store.Create("t", "d", state)
		time.Sleep(10 * time.Millisecond)
	}

	// Test MaxCount: 3
	err = store.Cleanup(RetentionConfig{MaxCount: 3, MaxAge: 24 * time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	snaps, _ := store.List()
	if len(snaps) != 3 {
		t.Errorf("Expected 3 snapshots, got %d", len(snaps))
	}

	// Test MaxAge: very short
	err = store.Cleanup(RetentionConfig{MaxCount: 10, MaxAge: 1 * time.Nanosecond})
	if err != nil {
		t.Fatal(err)
	}
	snaps, _ = store.List()
	if len(snaps) != 0 {
		t.Errorf("Expected 0 snapshots after age cleanup, got %d", len(snaps))
	}
}
