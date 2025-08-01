package paths

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	dataDir string
	once    sync.Once
)

// SetDataDir sets a custom data directory
// This should be called before any other paths functions
func SetDataDir(dir string) error {
	if dir == "" {
		return nil
	}
	
	// Expand ~ to home directory
	if len(dir) >= 2 && dir[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dir = filepath.Join(homeDir, dir[2:])
	}
	
	// Make absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(absDir, 0755); err != nil {
		return err
	}
	
	dataDir = absDir
	return nil
}

// initDataDir initializes the data directory if not already set
func initDataDir() {
	once.Do(func() {
		if dataDir != "" {
			return
		}
		
		// Priority order:
		// 1. Already set via SetDataDir
		// 2. Development default (if built with -tags dev)
		// 3. SHIEN_DATA_DIR environment variable
		// 4. Default ~/.config/shien
		
		// Try development default first
		if devDir := getDefaultDataDir(); devDir != "" {
			if err := SetDataDir(devDir); err == nil {
				return
			}
		}
		
		// Try environment variable
		if envDir := os.Getenv("SHIEN_DATA_DIR"); envDir != "" {
			if err := SetDataDir(envDir); err == nil {
				return
			}
		}
		
		// Default to ~/.config/shien
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		dataDir = filepath.Join(homeDir, ".config", "shien")
		
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			panic(err)
		}
	})
}

func DataDir() string {
	initDataDir()
	return dataDir
}

func ConfigFile() string {
	initDataDir()
	return filepath.Join(dataDir, "config.json")
}

func DatabaseFile() string {
	initDataDir()
	return filepath.Join(dataDir, "shien.db")
}

func SocketFile() string {
	initDataDir()
	return filepath.Join(dataDir, "shien-service.sock")
}

func IsDevMode() bool {
	initDataDir()
	// Dev mode if not using default path
	homeDir, _ := os.UserHomeDir()
	defaultPath := filepath.Join(homeDir, ".config", "shien")
	return dataDir != defaultPath
}