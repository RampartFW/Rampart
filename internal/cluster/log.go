package cluster

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/rampartfw/rampart/internal/model"
)

// Log manages the Raft log and its persistence on disk.
type Log struct {
	mu      sync.RWMutex
	entries []model.LogEntry
	wal     *os.File
	encoder *gob.Encoder
	path    string
}

// NewLog creates a new Raft log with persistence.
func NewLog(path string) (*Log, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAL file: %w", err)
	}

	l := &Log{
		entries: make([]model.LogEntry, 0),
		wal:     file,
		encoder: gob.NewEncoder(file),
		path:    path,
	}

	if err := l.load(); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to load WAL: %w", err)
	}

	return l, nil
}

// load reads the WAL from disk and restores the log in memory.
func (l *Log) load() error {
	_, err := l.wal.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(l.wal)
	for {
		var entry model.LogEntry
		if err := decoder.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		l.entries = append(l.entries, entry)
	}

	return nil
}

// Append adds new entries to the log and persists them.
func (l *Log) Append(entries ...model.LogEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, entry := range entries {
		if err := l.encoder.Encode(entry); err != nil {
			return fmt.Errorf("failed to persist entry: %w", err)
		}
		l.entries = append(l.entries, entry)
	}

	return nil
}

// Get returns an entry at a specific index.
func (l *Log) Get(index uint64) (model.LogEntry, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if index == 0 || index > uint64(len(l.entries)) {
		return model.LogEntry{}, false
	}
	return l.entries[index-1], true
}

// LastEntry returns the last entry in the log.
func (l *Log) LastEntry() model.LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(l.entries) == 0 {
		return model.LogEntry{}
	}
	return l.entries[len(l.entries)-1]
}

// Entries returns a slice of entries from start to end (inclusive, 1-based index).
func (l *Log) Entries(start, end uint64) []model.LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if start == 0 {
		start = 1
	}
	if end > uint64(len(l.entries)) {
		end = uint64(len(l.entries))
	}
	if start > end {
		return nil
	}

	result := make([]model.LogEntry, end-start+1)
	copy(result, l.entries[start-1:end])
	return result
}

// Truncate removes all entries from index onwards.
func (l *Log) Truncate(index uint64) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if index == 0 || index > uint64(len(l.entries)) {
		return nil
	}

	l.entries = l.entries[:index-1]

	// Rewrite WAL
	l.wal.Close()
	file, err := os.OpenFile(l.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	l.wal = file
	l.encoder = gob.NewEncoder(file)

	for _, entry := range l.entries {
		if err := l.encoder.Encode(entry); err != nil {
			return err
		}
	}

	return nil
}

// Close closes the WAL file.
func (l *Log) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.wal.Close()
}
