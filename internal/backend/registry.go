package backend

import (
	"fmt"
	"strings"
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

// NewBackend creates a new backend by name. 
// Supports comma-separated names for multi-backend fan-out (e.g. "nftables,aws").
func NewBackend(name string, cfg BackendConfig) (Backend, error) {
	if strings.Contains(name, ",") {
		names := strings.Split(name, ",")
		var backends []Backend
		for _, n := range names {
			n = strings.TrimSpace(n)
			b, err := NewBackend(n, cfg)
			if err != nil {
				return nil, fmt.Errorf("failed to create sub-backend %s: %w", n, err)
			}
			backends = append(backends, b)
		}
		return NewMultiBackend(backends...), nil
	}

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
	// 1. Try to initialize Hybrid (eBPF for fast path + nftables for slow path)
	nft, err1 := NewBackend("nftables", BackendConfig{Type: "nftables"})
	ebpf, err2 := NewBackend("ebpf", BackendConfig{Type: "ebpf"})

	if err1 == nil && err2 == nil && nft.Probe() == nil && ebpf.Probe() == nil {
		// Both supported! Create hybrid. 
		if factory, ok := registry["hybrid"]; ok {
			return factory(BackendConfig{Type: "hybrid"})
		}
	}

	// Priority fallback: nftables > iptables > ebpf
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
