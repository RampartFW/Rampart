package backend

import (
	"fmt"
	"sync"
)

// BackendFactory is a function that creates a new Backend
type BackendFactory func(cfg BackendConfig) (Backend, error)

var (
	registry = make(map[string]BackendFactory)
	mu       sync.RWMutex
)

// Register registers a new backend factory
func Register(name string, factory BackendFactory) {
	mu.Lock()
	defer mu.Unlock()
	registry[name] = factory
}

// NewBackend creates a new backend by name
func NewBackend(name string, cfg BackendConfig) (Backend, error) {
	mu.RLock()
	factory, ok := registry[name]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown backend: %s", name)
	}
	return factory(cfg)
}

// AutoDetect probes for the best available backend
func AutoDetect() (Backend, error) {
	// Priority: nftables > iptables > ebpf
	for _, name := range []string{"nftables", "iptables", "ebpf"} {
		mu.RLock()
		factory, ok := registry[name]
		mu.RUnlock()
		if !ok {
			continue
		}

		b, err := factory(BackendConfig{Type: name})
		if err == nil && b.Probe() == nil {
			return b, nil
		}
	}
	return nil, fmt.Errorf("no supported firewall backend found")
}
