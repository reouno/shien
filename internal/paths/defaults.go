//go:build !dev

package paths

// getDefaultDataDir returns the default data directory
// This is the production version
func getDefaultDataDir() string {
	return "" // Use system default
}