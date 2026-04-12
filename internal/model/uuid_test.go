package model

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateUUIDv7(t *testing.T) {
	uuid := GenerateUUIDv7()

	// Check format (8-4-4-4-12)
	parts := strings.Split(uuid, "-")
	if len(parts) != 5 {
		t.Fatalf("expected 5 parts, got %d", len(parts))
	}
	if len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
		t.Fatalf("invalid part lengths: %v", parts)
	}

	// Check version 7
	// uuid[6] is parts[2]
	if parts[2][0] != '7' {
		t.Errorf("expected version 7, got %c", parts[2][0])
	}

	// Check variant (bits 64-65)
	// parts[3] is bytes 8-9. First digit of parts[3] should be 8, 9, a, or b.
	variantDigit := parts[3][0]
	if variantDigit != '8' && variantDigit != '9' && variantDigit != 'a' && variantDigit != 'b' {
		t.Errorf("expected variant RFC 4122, got %c", variantDigit)
	}
}

func TestGenerateUUIDv7Monotonicity(t *testing.T) {
	u1 := GenerateUUIDv7()
	u2 := GenerateUUIDv7()

	if u1 >= u2 {
		t.Errorf("expected u2 > u1 for monotonicity, got u1=%s, u2=%s", u1, u2)
	}

	// Test across a millisecond boundary
	time.Sleep(2 * time.Millisecond)
	u3 := GenerateUUIDv7()
	if u2 >= u3 {
		t.Errorf("expected u3 > u2 after sleep, got u2=%s, u3=%s", u2, u3)
	}
}
