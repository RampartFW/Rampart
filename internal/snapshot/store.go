package snapshot

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/rampartfw/rampart/internal/engine"
	"github.com/rampartfw/rampart/internal/model"
)

// Store manages snapshot storage and retrieval.
type Store struct {
	dir string
	mu  sync.RWMutex
}

// NewStore creates a new snapshot store.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create snapshot directory: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Create captures a new snapshot of the current ruleset.
func (s *Store) Create(trigger, description string, state *model.CompiledRuleSet) (*model.Snapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := model.GenerateUUIDv7()
	filename := fmt.Sprintf("%s-%s.snap", id, trigger)
	path := filepath.Join(s.dir, filename)

	// Save the ruleset using Gob encoding
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot file: %w", err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(state); err != nil {
		return nil, fmt.Errorf("failed to encode ruleset: %w", err)
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat snapshot file: %w", err)
	}

	snap := &model.Snapshot{
		ID:          id,
		CreatedAt:   time.Now(),
		CreatedBy:   "system", // In a real app, this would be passed from context
		Trigger:     trigger,
		Description: description,
		PolicyHash:  state.Hash,
		RuleCount:   len(state.Rules),
		Backend:     state.Backend,
		Size:        fi.Size(),
		Filename:    filename,
	}

	// Update the index
	if err := s.addToIndex(snap); err != nil {
		return nil, fmt.Errorf("failed to update snapshot index: %w", err)
	}

	return snap, nil
}

// List returns all snapshots, sorted by creation time (newest first).
func (s *Store) List() ([]model.Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.readIndex()
}

// Load retrieves a snapshot and its associated ruleset.
func (s *Store) Load(id string) (*model.Snapshot, *model.CompiledRuleSet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snaps, err := s.readIndex()
	if err != nil {
		return nil, nil, err
	}

	var snap *model.Snapshot
	for i := range snaps {
		if snaps[i].ID == id {
			snap = &snaps[i]
			break
		}
	}

	if snap == nil {
		return nil, nil, fmt.Errorf("snapshot not found: %s", id)
	}

	path := filepath.Join(s.dir, snap.Filename)
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open snapshot file: %w", err)
	}
	defer f.Close()

	var state model.CompiledRuleSet
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&state); err != nil {
		return nil, nil, fmt.Errorf("failed to decode ruleset: %w", err)
	}

	return snap, &state, nil
}

// Delete removes a snapshot and its associated file.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	snaps, err := s.readIndex()
	if err != nil {
		return err
	}

	newSnaps := make([]model.Snapshot, 0, len(snaps))
	var snapToDelete *model.Snapshot

	for i := range snaps {
		if snaps[i].ID == id {
			snapToDelete = &snaps[i]
		} else {
			newSnaps = append(newSnaps, snaps[i])
		}
	}

	if snapToDelete == nil {
		return fmt.Errorf("snapshot not found: %s", id)
	}

	// Delete file
	path := filepath.Join(s.dir, snapToDelete.Filename)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete snapshot file: %w", err)
	}

	// Update index
	return s.writeIndex(newSnaps)
}

// Diff returns an ExecutionPlan showing the differences between a snapshot and a ruleset.
func (s *Store) Diff(id string, current *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	_, snapState, err := s.Load(id)
	if err != nil {
		return nil, err
	}

	return engine.GeneratePlan(snapState, current), nil
}

// Index management helpers

func (s *Store) addToIndex(snap *model.Snapshot) error {
	snaps, err := s.readIndex()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	snaps = append(snaps, *snap)
	return s.writeIndex(snaps)
}

func (s *Store) readIndex() ([]model.Snapshot, error) {
	path := filepath.Join(s.dir, "snapshots.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Snapshot{}, nil
		}
		return nil, err
	}

	var snaps []model.Snapshot
	if err := json.Unmarshal(data, &snaps); err != nil {
		return nil, err
	}

	// Sort by creation time (newest first)
	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].CreatedAt.After(snaps[j].CreatedAt)
	})

	return snaps, nil
}

func (s *Store) writeIndex(snaps []model.Snapshot) error {
	path := filepath.Join(s.dir, "snapshots.json")
	data, err := json.MarshalIndent(snaps, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
