package service

import (
	"shien/internal/config"
	"shien/internal/database"
)

// Services aggregates all service layers
type Services struct {
	Activity *ActivityService
	Config   *ConfigService
}

// NewServices creates all services
func NewServices(repo *database.Repository, cfg *config.Manager) *Services {
	return &Services{
		Activity: NewActivityService(repo.Activity()),
		Config:   NewConfigService(cfg),
	}
}

// ConfigService handles configuration logic
type ConfigService struct {
	manager *config.Manager
}

// NewConfigService creates a new config service
func NewConfigService(manager *config.Manager) *ConfigService {
	return &ConfigService{manager: manager}
}

// GetConfig returns current configuration
func (s *ConfigService) GetConfig() *config.Config {
	if s.manager == nil {
		return config.DefaultConfig()
	}
	return s.manager.Get()
}

// UpdateConfig updates configuration
func (s *ConfigService) UpdateConfig(updates map[string]interface{}) error {
	// Future implementation
	return nil
}