//go:build linux

package ebpf

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

type EBPFBackend struct {
	cfg backend.BackendConfig
	iface string
	progFd int
	maps map[string]int
}

func init() {
	backend.Register("ebpf", func(cfg backend.BackendConfig) (backend.Backend, error) {
		iface, _ := cfg.Settings["interface"]
		return &EBPFBackend{
			cfg: cfg,
			iface: iface,
			maps: make(map[string]int),
		}, nil
	})
}

func (b *EBPFBackend) Name() string {
	return "ebpf"
}

func (b *EBPFBackend) Capabilities() model.BackendCapabilities {
	return model.BackendCapabilities{
		IPv4:               true,
		IPv6:               true,
		RateLimiting:       true,
		ConnectionTracking: false, // Limited in eBPF for this MVP
		Logging:            true,
		PerRuleCounters:    true,
		AtomicReplace:      false,
		InterfaceFiltering: true,
	}
}

func (b *EBPFBackend) Probe() error {
	// Check if we can call bpf()
	_, _, err := unix.Syscall(unix.SYS_BPF, 0, 0, 0)
	if err != 0 && err != unix.EPERM && err != unix.EINVAL {
		return fmt.Errorf("ebpf kernel support check failed: %w", err)
	}
	return nil
}

func (b *EBPFBackend) CurrentState() (*model.CompiledRuleSet, error) {
	// eBPF state is in maps. It's hard to reconstruct exactly.
	// We'll return empty for now, or just the metadata.
	return &model.CompiledRuleSet{
		Rules: []model.CompiledRule{},
	}, nil
}

func (b *EBPFBackend) Apply(rs *model.CompiledRuleSet) error {
	if b.iface == "" {
		return fmt.Errorf("interface must be specified for ebpf backend")
	}

	// 1. Create maps if not exist
	if err := b.createMaps(); err != nil {
		return err
	}

	// 2. Load program
	if err := b.loadProgram(); err != nil {
		return err
	}

	// 3. Attach to interface
	if err := b.attachXDP(); err != nil {
		return err
	}

	// 4. Update maps with rules
	return b.updateMaps(rs)
}

func (b *EBPFBackend) createMaps() error {
	// LPM Trie for CIDRs
	fd, err := b.createMap(unix.BPF_MAP_TYPE_LPM_TRIE, 8, 4, 1024, unix.BPF_F_NO_PREALLOC)
	if err != nil {
		return fmt.Errorf("failed to create lpm_trie map: %w", err)
	}
	b.maps["blocked_cidrs"] = fd

	// Hash map for ports
	fd, err = b.createMap(unix.BPF_MAP_TYPE_HASH, 2, 4, 1024, 0)
	if err != nil {
		return fmt.Errorf("failed to create ports map: %w", err)
	}
	b.maps["allowed_ports"] = fd

	// Percpu array for stats
	fd, err = b.createMap(unix.BPF_MAP_TYPE_PERCPU_ARRAY, 4, 16, 1024, 0)
	if err != nil {
		return fmt.Errorf("failed to create stats map: %w", err)
	}
	b.maps["rule_stats"] = fd

	return nil
}

func (b *EBPFBackend) createMap(mapType int, keySize, valueSize, maxEntries, flags uint32) (int, error) {
	attr := struct {
		mapType    uint32
		keySize    uint32
		valueSize  uint32
		maxEntries uint32
		mapFlags   uint32
	}{
		mapType:    mapType,
		keySize:    keySize,
		valueSize:  valueSize,
		maxEntries: maxEntries,
		mapFlags:   flags,
	}

	fd, _, err := unix.Syscall(unix.SYS_BPF, unix.BPF_MAP_CREATE, uintptr(unsafe.Pointer(&attr)), unsafe.Sizeof(attr))
	if err != 0 {
		return 0, err
	}
	return int(fd), nil
}

func (b *EBPFBackend) loadProgram() error {
	// Placeholder bytecode logic
	bytecode := getXDPBytecode() // Defined in xdp_program.go

	attr := struct {
		progType    uint32
		insnCnt     uint32
		insns       uintptr
		license     uintptr
		logLevel    uint32
		logSize     uint32
		logBuf      uintptr
		kernVersion uint32
		progFlags   uint32
	}{
		progType: unix.BPF_PROG_TYPE_XDP,
		insnCnt:  uint32(len(bytecode) / 8),
		insns:    uintptr(unsafe.Pointer(&bytecode[0])),
		license:  uintptr(unsafe.Pointer(&[]byte("GPL\x00")[0])),
	}

	fd, _, err := unix.Syscall(unix.SYS_BPF, unix.BPF_PROG_LOAD, uintptr(unsafe.Pointer(&attr)), unsafe.Sizeof(attr))
	if err != 0 {
		return err
	}
	b.progFd = int(fd)
	return nil
}

func (b *EBPFBackend) attachXDP() error {
	// Simplified XDP attachment
	// In a real implementation, we'd use netlink to set IFLA_XDP_FD
	return nil
}

func (b *EBPFBackend) updateMaps(rs *model.CompiledRuleSet) error {
	// Implement rule translation to map updates
	return nil
}

func (b *EBPFBackend) DryRun(rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return &model.ExecutionPlan{
		PlannedRuleCount: len(rs.Rules),
	}, nil
}

func (b *EBPFBackend) Rollback(snapshot *model.Snapshot) error {
	return fmt.Errorf("rollback not implemented for ebpf")
}

func (b *EBPFBackend) Flush() error {
	// Close FDs
	for _, fd := range b.maps {
		unix.Close(fd)
	}
	if b.progFd != 0 {
		unix.Close(b.progFd)
	}
	return nil
}

func (b *EBPFBackend) Stats() (map[string]model.RuleStats, error) {
	return make(map[string]model.RuleStats), nil
}

func (b *EBPFBackend) Close() error {
	return b.Flush()
}
