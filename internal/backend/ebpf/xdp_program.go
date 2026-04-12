//go:build linux

package ebpf

// getXDPBytecode returns a placeholder XDP program bytecode.
// In a real implementation, this would be either pre-compiled ELF bytes
// or bytecode generated at runtime.
func getXDPBytecode() []byte {
	// A very simple XDP program that just returns XDP_PASS (2)
	// Bytecode for: return XDP_PASS;
	// 64-bit BPF instructions
	return []byte{
		0xb7, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, // r0 = 2 (XDP_PASS)
		0x95, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // exit
	}
}
