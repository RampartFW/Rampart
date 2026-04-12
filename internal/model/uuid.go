package model

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

var (
	lastMs   uint64
	lastSeq  uint16
	uuidMu   sync.Mutex
)

// GenerateUUIDv7 generates a new UUID v7 (time-sortable).
// It uses a 48-bit timestamp (milliseconds), 12-bit sequence for monotonicity within the same ms,
// and 62 bits of randomness.
func GenerateUUIDv7() string {
	uuidMu.Lock()
	defer uuidMu.Unlock()

	now := time.Now().UnixMilli()
	ms := uint64(now)

	if ms <= lastMs {
		ms = lastMs
		lastSeq++
	} else {
		lastMs = ms
		lastSeq = 0
	}

	var uuid [16]byte

	// 48-bit timestamp
	uuid[0] = byte(ms >> 40)
	uuid[1] = byte(ms >> 32)
	uuid[2] = byte(ms >> 24)
	uuid[3] = byte(ms >> 16)
	uuid[4] = byte(ms >> 8)
	uuid[5] = byte(ms)

	// Version 7 (bits 48-51)
	uuid[6] = 0x70 | (byte(lastSeq>>8) & 0x0F)
	uuid[7] = byte(lastSeq)

	// Variant 2 (RFC 4122) (bits 64-65)
	// Random suffix (remaining 62 bits)
	_, _ = rand.Read(uuid[8:])
	uuid[8] = 0x80 | (uuid[8] & 0x3F)

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}
