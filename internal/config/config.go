package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Config holds application configuration
type Config struct {
	// Notification settings
	NotificationEnabled   bool   `json:"notification_enabled"`
	NotificationSound     string `json:"notification_sound"`
	
	// Application settings
	StartOnLogin          bool   `json:"start_on_login"`
	ShowInDock            bool   `json:"show_in_dock"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		NotificationEnabled:   true,
		NotificationSound:     "default",
		StartOnLogin:          false,
		ShowInDock:            false,
	}
}

// Manager handles configuration persistence
type Manager struct {
	config     *Config
	configPath string
	mu         sync.RWMutex
}

// NewManager creates a new config manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	configDir := filepath.Join(homeDir, ".config", "shien")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}
	
	configPath := filepath.Join(configDir, "config.json")
	
	m := &Manager{
		configPath: configPath,
		config:     DefaultConfig(),
	}
	
	// Load existing config if available
	if err := m.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	
	// Save default config if not exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := m.Save(); err != nil {
			return nil, err
		}
	}
	
	return m, nil
}

// Get returns current configuration
func (m *Manager) Get() *Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Return a copy to prevent external modifications
	cfg := *m.config
	return &cfg
}

// Update updates configuration
func (m *Manager) Update(fn func(*Config)) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	fn(m.config)
	return m.Save()
}

// Load loads configuration from file
func (m *Manager) Load() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, m.config)
}

// Save saves configuration to file
func (m *Manager) Save() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(m.configPath, data, 0644)
}

// ConfigPath returns the path to config file
func (m *Manager) ConfigPath() string {
	return m.configPath
}