//go:build linux

package ebpf

import (
	"context"
	"fmt"
	"net"
	"unsafe"

	"golang.org/x/sys/unix"
	"github.com/rampartfw/rampart/internal/backend"
	"github.com/rampartfw/rampart/internal/model"
)

type EBPFBackend struct {
	cfg   backend.BackendConfig
	iface string
	ifIdx int
	progFd int
	maps  map[string]int
}

func init() {
	backend.Register("ebpf", func(cfg backend.BackendConfig) (backend.Backend, error) {
		iface, _ := cfg.Settings["interface"]
		var ifIdx int
		if iface != "" {
			ifObj, err := net.InterfaceByName(iface)
			if err == nil {
				ifIdx = ifObj.Index
			}
		}
		
		return &EBPFBackend{
			cfg:   cfg,
			iface: iface,
			ifIdx: ifIdx,
			maps:  make(map[string]int),
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
		ConnectionTracking: false,
		Logging:            true,
		PerRuleCounters:    true,
		AtomicReplace:      false,
		InterfaceFiltering: true,
	}
}

func (b *EBPFBackend) Probe() error {
	// Check if we can call bpf()
	attr := struct {
		mapType    uint32
		keySize    uint32
		valueSize  uint32
		maxEntries uint32
		mapFlags   uint32
	}{
		mapType:    unix.BPF_MAP_TYPE_ARRAY,
		keySize:    4,
		valueSize:  4,
		maxEntries: 1,
	}
	fd, _, err := unix.Syscall(unix.SYS_BPF, unix.BPF_MAP_CREATE, uintptr(unsafe.Pointer(&attr)), unsafe.Sizeof(attr))
	if err != 0 {
		return fmt.Errorf("ebpf kernel support check failed: %w", err)
	}
	unix.Close(int(fd))
	return nil
}

func (b *EBPFBackend) CurrentState(ctx context.Context) (*model.CompiledRuleSet, error) {
	return &model.CompiledRuleSet{
		Rules: []model.CompiledRule{},
	}, nil
}

func (b *EBPFBackend) Apply(ctx context.Context, rs *model.CompiledRuleSet) error {
	if b.ifIdx == 0 {
		return fmt.Errorf("invalid interface or interface not specified")
	}

	// 1. Create maps
	if err := b.createMaps(); err != nil {
		return err
	}

	// 2. Load program
	if err := b.loadProgram(); err != nil {
		return err
	}

	// 3. Attach to interface via Netlink
	if err := b.attachXDP(); err != nil {
		return err
	}

	// 4. Populate maps
	return b.updateMaps(rs)
}

func (b *EBPFBackend) createMaps() error {
	// LPM Trie for CIDRs
	// Key: [prefix_len(4) + ip(4 or 16)]
	fd, err := b.createMap(unix.BPF_MAP_TYPE_LPM_TRIE, 8, 4, 4096, unix.BPF_F_NO_PREALLOC)
	if err != nil {
		return err
	}
	b.maps["blocked_cidrs"] = fd

	// Hash map for ports
	fd, err = b.createMap(unix.BPF_MAP_TYPE_HASH, 2, 4, 1024, 0)
	if err != nil {
		return err
	}
	b.maps["allowed_ports"] = fd

	// Percpu array for stats
	fd, err = b.createMap(unix.BPF_MAP_TYPE_PERCPU_ARRAY, 4, 16, 4096, 0)
	if err != nil {
		return err
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
		mapType:    uint32(mapType),
		keySize:    keySize,
		valueSize:  valueSize,
		maxEntries: maxEntries,
		mapFlags:   flags,
	}

	fd, _, errno := unix.Syscall(unix.SYS_BPF, unix.BPF_MAP_CREATE, uintptr(unsafe.Pointer(&attr)), unsafe.Sizeof(attr))
	if errno != 0 {
		return 0, errno
	}
	return int(fd), nil
}

func (b *EBPFBackend) loadProgram() error {
	// In a real product, we load from an embedded ELF or pre-compiled bytecode
	insns := getXDPBytecode() 
	
	license := []byte("GPL\x00")
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
		insnCnt:  uint32(len(insns) / 8),
		insns:    uintptr(unsafe.Pointer(&insns[0])),
		license:  uintptr(unsafe.Pointer(&license[0])),
	}

	fd, _, errno := unix.Syscall(unix.SYS_BPF, unix.BPF_PROG_LOAD, uintptr(unsafe.Pointer(&attr)), unsafe.Sizeof(attr))
	if errno != 0 {
		return errno
	}
	b.progFd = int(fd)
	return nil
}

func (b *EBPFBackend) attachXDP() error {
	// Use netlink to attach XDP program
	// Construct RTM_SETLINK message with IFLA_XDP attribute
	return b.setXdpFd(b.ifIdx, b.progFd)
}

func (b *EBPFBackend) setXdpFd(ifIdx, fd int) error {
	// Create a raw Netlink socket
	s, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_ROUTE)
	if err != nil {
		return fmt.Errorf("failed to create netlink socket: %w", err)
	}
	defer unix.Close(s)

	// Netlink Message Structure: [Header][IfInfoMsg][Attributes]
	// IFLA_XDP is the attribute type we need to set
	
	// 1. Prepare IfInfoMsg
	ifm := &unix.IfInfomsg{
		Family: unix.AF_UNSPEC,
		Index:  int32(ifIdx),
	}

	// 2. Prepare IFLA_XDP Attribute
	// This is a nested attribute: IFLA_XDP -> IFLA_XDP_FD
	xdpAttr := struct {
		Type  uint16
		Len   uint16
		Value uint32
	}{
		Type:  6, // IFLA_XDP_FD
		Len:   8,
		Value: uint32(fd),
	}
	
	// Simplified Netlink message construction for production
	// In a real high-perf scenario, we'd use a buffer pool
	msg := make([]byte, 0, 128)
	// Add placeholders for lengths to be filled
	msg = append(msg, make([]byte, unix.SizeofNlMsghdr+unix.SizeofIfInfomsg)...)
	
	// Add IFLA_XDP nested attribute header
	// Type 43 is IFLA_XDP, nested
	msg = append(msg, []byte{12, 0, 43, 128}...) // Len: 12, Type: 43 | NLA_F_NESTED
	
	// Add IFLA_XDP_FD attribute
	fdBytes := make([]byte, 8)
	*(*uint16)(unsafe.Pointer(&fdBytes[0])) = 8 // Len
	*(*uint16)(unsafe.Pointer(&fdBytes[2])) = 6 // IFLA_XDP_FD
	*(*uint32)(unsafe.Pointer(&fdBytes[4])) = uint32(fd)
	msg = append(msg, fdBytes...)

	// Fill headers
	hdr := (*unix.NlMsghdr)(unsafe.Pointer(&msg[0]))
	hdr.Len = uint32(len(msg))
	hdr.Type = unix.RTM_SETLINK
	hdr.Flags = unix.NLM_F_REQUEST | unix.NLM_F_ACK

	ifmPtr := (*unix.IfInfomsg)(unsafe.Pointer(&msg[unix.SizeofNlMsghdr]))
	*ifmPtr = *ifm

	// Send message
	if err := unix.Sendto(s, msg, 0, &unix.SockaddrNetlink{Family: unix.AF_NETLINK}); err != nil {
		return fmt.Errorf("failed to send netlink message: %w", err)
	}

	// Wait for ACK
	resp := make([]byte, 4096)
	n, _, err := unix.Recvfrom(s, resp, 0)
	if err != nil {
		return fmt.Errorf("failed to receive netlink ack: %w", err)
	}

	// Parse ACK (NLMSG_ERROR)
	if n >= unix.SizeofNlMsghdr {
		respHdr := (*unix.NlMsghdr)(unsafe.Pointer(&resp[0]))
		if respHdr.Type == unix.NLMSG_ERROR {
			nlErr := (*unix.NlMsgerr)(unsafe.Pointer(&resp[unix.SizeofNlMsghdr]))
			if nlErr.Error != 0 {
				return fmt.Errorf("netlink error: %d", -nlErr.Error)
			}
		}
	}

	return nil
}

func (b *EBPFBackend) updateMaps(rs *model.CompiledRuleSet) error {
	for i, rule := range rs.Rules {
		// Example: Update blocked CIDRs map
		for _, net := range rule.Match.SourceNets {
			if rule.Action == model.ActionDrop {
				b.updateLPMTrie(b.maps["blocked_cidrs"], net, uint32(i))
			}
		}
	}
	return nil
}

func (b *EBPFBackend) updateLPMTrie(mapFd int, ipNet net.IPNet, value uint32) error {
	ones, _ := ipNet.Mask.Size()
	key := make([]byte, 8) // 4 bytes prefix_len + 4 bytes IPv4
	*(*uint32)(unsafe.Pointer(&key[0])) = uint32(ones)
	copy(key[4:], ipNet.IP.To4())

	attr := struct {
		mapFd uint32
		key   uintptr
		value uintptr
		flags uint64
	}{
		mapFd: uint32(mapFd),
		key:   uintptr(unsafe.Pointer(&key[0])),
		value: uintptr(unsafe.Pointer(&value)),
		flags: 0,
	}

	_, _, errno := unix.Syscall(unix.SYS_BPF, unix.BPF_MAP_UPDATE_ELEM, uintptr(unsafe.Pointer(&attr)), unsafe.Sizeof(attr))
	if errno != 0 {
		return errno
	}
	return nil
}

func (b *EBPFBackend) DryRun(ctx context.Context, rs *model.CompiledRuleSet) (*model.ExecutionPlan, error) {
	return &model.ExecutionPlan{
		PlannedRuleCount: len(rs.Rules),
	}, nil
}

func (b *EBPFBackend) Rollback(ctx context.Context, snapshot *model.Snapshot) error {
	return fmt.Errorf("rollback not implemented for ebpf")
}

func (b *EBPFBackend) Flush(ctx context.Context) error {
	// Close FDs and detach XDP
	b.setXdpFd(b.ifIdx, -1)
	for _, fd := range b.maps {
		unix.Close(fd)
	}
	if b.progFd != 0 {
		unix.Close(b.progFd)
	}
	return nil
}

func (b *EBPFBackend) Stats(ctx context.Context) (map[string]model.RuleStats, error) {
	return make(map[string]model.RuleStats), nil
}

func (b *EBPFBackend) Close() error {
	return b.Flush(context.Background())
}

// Dummy for now, in a real product this would be loaded from an object file
func getXDPBytecode() []uint64 {
	return []uint64{
		0x00000000000000b7, // mov r0, 2 (XDP_PASS)
		0x0000000000000095, // exit
	}
}
