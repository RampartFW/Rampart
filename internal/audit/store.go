package audit

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rampartfw/rampart/internal/model"
)

// Store manages audit log storage and retrieval.
type Store struct {
	dir      string
	mu       sync.RWMutex
	eventC   chan model.AuditEvent
	lastHash string
	maxAge   time.Duration
	done     chan struct{}
	wg       sync.WaitGroup
}

// NewStore creates a new audit store.
func NewStore(dir string, maxAge time.Duration) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create audit directory: %w", err)
	}

	s := &Store{
		dir:    dir,
		eventC: make(chan model.AuditEvent, 1000),
		maxAge: maxAge,
		done:   make(chan struct{}),
	}

	// Initialize lastHash from the latest file
	if err := s.initLastHash(); err != nil {
		return nil, err
	}

	s.wg.Add(2)
	go s.runWriter()
	go s.runMaintenance()

	return s, nil
}

// Record appends an audit event to the log.
func (s *Store) Record(event model.AuditEvent) error {
	if event.ID == "" {
		event.ID = model.GenerateUUIDv7()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	select {
	case s.eventC <- event:
		return nil
	default:
		return fmt.Errorf("audit event channel full")
	}
}

// Get retrieves a specific audit event by ID.
func (s *Store) Get(id string) (*model.AuditEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	files, err := s.listAuditFiles()
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		event, err := s.findInFile(f, id)
		if err == nil {
			return event, nil
		}
	}

	return nil, fmt.Errorf("audit event not found: %s", id)
}

// Search filters audit events based on the query.
func (s *Store) Search(query AuditQuery) ([]model.AuditEvent, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []model.AuditEvent
	files, err := s.listAuditFiles()
	if err != nil {
		return nil, 0, err
	}

	for _, f := range files {
		if !s.fileMightContainRange(f, query.Since, query.Until) {
			continue
		}

		events, err := s.readAllFromFile(f)
		if err != nil {
			log.Printf("Failed to read audit file %s: %v", f, err)
			continue
		}

		for _, e := range events {
			if query.Matches(e) {
				results = append(results, e)
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp.After(results[j].Timestamp)
	})

	total := len(results)
	start := query.Offset
	if start > total {
		start = total
	}
	end := start + query.Limit
	if query.Limit == 0 || end > total {
		end = total
	}

	return results[start:end], total, nil
}

// VerifyIntegrity verifies the hash chain across all audit logs.
func (s *Store) VerifyIntegrity() (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	files, err := s.listAuditFiles()
	if err != nil {
		return false, err
	}

	sort.Strings(files)

	currentHash := ""
	for _, f := range files {
		events, err := s.readAllFromFile(f)
		if err != nil {
			return false, err
		}

		for _, e := range events {
			hash := e.ChainHash
			e.ChainHash = ""
			entryJSON, _ := json.Marshal(e)
			
			computedHash := s.computeHash(currentHash, entryJSON)
			if computedHash != hash {
				return false, fmt.Errorf("integrity violation at event %s in file %s", e.ID, f)
			}
			currentHash = hash
		}
	}

	return true, nil
}

func (s *Store) runWriter() {
	defer s.wg.Done()
	for event := range s.eventC {
		s.mu.Lock()
		if err := s.writeEvent(event); err != nil {
			log.Printf("Failed to write audit event: %v", err)
		}
		s.mu.Unlock()
	}
}

func (s *Store) writeEvent(event model.AuditEvent) (err error) {
	today := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("audit-%s.jsonl", today)
	path := filepath.Join(s.dir, filename)

	event.ChainHash = ""
	entryJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	event.ChainHash = s.computeHash(s.lastHash, entryJSON)
	s.lastHash = event.ChainHash

	finalJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if _, err = f.Write(finalJSON); err != nil {
		return err
	}
	if _, err = f.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

func (s *Store) computeHash(prevHash string, entryJSON []byte) string {
	h := sha256.New()
	h.Write([]byte(prevHash))
	h.Write(entryJSON)
	return hex.EncodeToString(h.Sum(nil))
}
func (s *Store) initLastHash() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	files, err := s.listAuditFiles()
	if err != nil {
		return err
	}

	if len(files) == 0 {
		s.lastHash = ""
		return nil
	}

	// Latest file
	latestFile := files[0]
	events, err := s.readAllFromFile(latestFile)
	if err != nil {
		return err
	}

	if len(events) > 0 {
		s.lastHash = events[len(events)-1].ChainHash
	}

	return nil
}

func (s *Store) listAuditFiles() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasPrefix(entry.Name(), "audit-") && (strings.HasSuffix(entry.Name(), ".jsonl") || strings.HasSuffix(entry.Name(), ".jsonl.gz"))) {
			files = append(files, filepath.Join(s.dir, entry.Name()))
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})

	return files, nil
}

func (s *Store) findInFile(path string, id string) (*model.AuditEvent, error) {
	events, err := s.readAllFromFile(path)
	if err != nil {
		return nil, err
	}

	for _, e := range events {
		if e.ID == id {
			return &e, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (s *Store) readAllFromFile(path string) ([]model.AuditEvent, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var r io.Reader = f
	if strings.HasSuffix(path, ".gz") {
		gzr, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		defer gzr.Close()
		r = gzr
	}

	var events []model.AuditEvent
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		var e model.AuditEvent
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			continue
		}
		events = append(events, e)
	}

	return events, scanner.Err()
}

func (s *Store) fileMightContainRange(path string, since, until time.Time) bool {
	base := filepath.Base(path)
	if !strings.HasPrefix(base, "audit-") {
		return true
	}
	dateStr := base[6:16]
	fileDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return true
	}

	if !since.IsZero() && fileDate.Add(24*time.Hour).Before(since) {
		return false
	}
	if !until.IsZero() && fileDate.After(until) {
		return false
	}
	return true
}

func (s *Store) runMaintenance() {
	defer s.wg.Done()
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.maintenance()
		case <-s.done:
			return
		}
	}
}

func (s *Store) maintenance() {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return
	}

	today := time.Now().Format("2006-01-02")
	cutoff := time.Now().Add(-s.maxAge)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, "audit-") {
			continue
		}

		dateStr := name[6:16]
		fileDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		if fileDate.Before(cutoff) {
			os.Remove(filepath.Join(s.dir, name))
			continue
		}

		if dateStr != today && !strings.HasSuffix(name, ".gz") {
			s.compressFile(filepath.Join(s.dir, name))
		}
	}
}

func (s *Store) compressFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	gzPath := path + ".gz"
	gzf, err := os.Create(gzPath)
	if err != nil {
		return
	}
	defer gzf.Close()

	gzw := gzip.NewWriter(gzf)
	if _, err := io.Copy(gzw, f); err != nil {
		gzw.Close()
		gzf.Close()
		os.Remove(gzPath)
		return
	}

	if err := gzw.Close(); err != nil {
		gzf.Close()
		os.Remove(gzPath)
		return
	}
	if err := gzf.Close(); err != nil {
		os.Remove(gzPath)
		return
	}
	f.Close()
	os.Remove(path)
}

// Close closes the store and waits for the writer to finish.
func (s *Store) Close() {
	close(s.eventC)
	close(s.done)
	s.wg.Wait()
}
