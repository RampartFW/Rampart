//go:build !linux

package security

func DropCapabilities() error {
	// Capability dropping is only implemented on Linux.
	return nil
}
