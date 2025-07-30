//go:build dev

package paths

import (
	"os"
	"path/filepath"
)

// getDefaultDataDir returns the default data directory
// This is the development version
func getDefaultDataDir() string {
	// Get executable path to find project root
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	
	// Assume we're in project root when building
	projectRoot := filepath.Dir(exePath)
	return filepath.Join(projectRoot, ".dev", "data")
}