//go:build !windows

package windowicon

// Apply is a no-op on platforms without Windows executable resources.
func Apply() {}
