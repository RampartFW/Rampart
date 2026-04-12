package snapshot

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rampartfw/rampart/internal/model"
)

// RetentionConfig defines how long snapshots are kept.
type RetentionConfig struct {
	MaxCount int
	MaxAge   time.Duration
}

// Cleanup deletes snapshots that no longer satisfy the retention policy.
// It keeps snapshots that are within BOTH the MaxCount and MaxAge limits.
func (s *Store) Cleanup(config RetentionConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	snaps, err := s.readIndex()
	if err != nil {
		return err
	}

	if len(snaps) == 0 {
		return nil
	}

	cutoff := time.Now().Add(-config.MaxAge)
	newSnaps := make([]model.Snapshot, 0, len(snaps))
	
	// Index is already sorted newest first in readIndex()
	// But let's be explicit here to be safe
	// Actually readIndex() does sort it.

	for i := range snaps {
		// Keep if it's within the count AND within the age
		if i < config.MaxCount && snaps[i].CreatedAt.After(cutoff) {
			newSnaps = append(newSnaps, snaps[i])
		} else {
			// Delete file
			path := filepath.Join(s.dir, snaps[i].Filename)
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				// We log but continue to delete other files and update index
				log.Printf("Failed to delete snapshot file %s: %v", path, err)
			}
		}
	}

	if len(newSnaps) == len(snaps) {
		return nil
	}

	return s.writeIndex(newSnaps)
}

// StartCleanupWorker starts a background goroutine to periodically clean up old snapshots.
func (s *Store) StartCleanupWorker(ctx context.Context, interval time.Duration, config RetentionConfig) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.Cleanup(config); err != nil {
				log.Printf("Snapshot cleanup failed: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
