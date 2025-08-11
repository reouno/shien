package service

import (
	"fmt"
	"shien/internal/database"
	"shien/internal/models/gamification"
	"time"

	"github.com/google/uuid"
)

// GamificationService handles gamification business logic
type GamificationService struct {
	repo   *database.Repository
	config *gamification.StatusConfig
}

// NewGamificationService creates a new gamification service
func NewGamificationService(repo *database.Repository) *GamificationService {
	return &GamificationService{
		repo:   repo,
		config: gamification.DefaultStatusConfig(),
	}
}

// GetConfig returns the gamification configuration
func (s *GamificationService) GetConfig() *gamification.StatusConfig {
	return s.config
}

// GetModifiers returns active modifiers for a user
func (s *GamificationService) GetModifiers(userID string) ([]gamification.AttributeModifier, error) {
	return s.repo.Gamification().GetAttributeModifiers(userID)
}

// GetOrCreateUserStatus retrieves or creates a user's status
func (s *GamificationService) GetOrCreateUserStatus(userID string) (*gamification.UserStatus, error) {
	// Try to get existing status
	status, err := s.repo.Gamification().GetUserStatus(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user status: %w", err)
	}
	
	// If status exists, return it
	if status != nil {
		return status, nil
	}
	
	// Create new status with default values
	newStatus := &gamification.UserStatus{
		UserID:        userID,
		Level:         1,
		Experience:    0,
		TotalExp:      0,
		Focus:         50,
		Productivity:  50,
		Creativity:    50,
		Stamina:       100,
		Knowledge:     10,
		Collaboration: 30,
	}
	
	err = s.repo.Gamification().CreateUserStatus(newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to create user status: %w", err)
	}
	
	return newStatus, nil
}

// ProcessActivity updates user status based on activity
func (s *GamificationService) ProcessActivity(userID string, appName string, duration time.Duration) error {
	// Get current status
	status, err := s.GetOrCreateUserStatus(userID)
	if err != nil {
		return err
	}
	
	// Get activity impact configuration
	impacts := gamification.PredefinedActivityImpacts()
	impact, exists := impacts[appName]
	if !exists {
		// Default impact for unknown apps
		impact = gamification.ActivityImpact{
			AppName:       appName,
			Category:      "other",
			FocusImpact:   0,
			ProductivityImpact: 1,
			CreativityImpact: 0,
			StaminaCost:   1,
			KnowledgeGain: 0,
			CollaborationImpact: 0,
			ExpGain:       3,
		}
	}
	
	// Calculate multiplier based on duration (5 minutes = 1x, 10 minutes = 2x, etc.)
	multiplier := int(duration.Minutes() / 5)
	if multiplier < 1 {
		multiplier = 1
	}
	
	// Apply impacts with multiplier
	status.Focus = gamification.ClampAttribute(status.Focus + (impact.FocusImpact * multiplier))
	status.Productivity = gamification.ClampAttribute(status.Productivity + (impact.ProductivityImpact * multiplier))
	status.Creativity = gamification.ClampAttribute(status.Creativity + (impact.CreativityImpact * multiplier))
	status.Stamina = gamification.ClampAttribute(status.Stamina - (impact.StaminaCost * multiplier))
	status.Knowledge = gamification.ClampAttribute(status.Knowledge + (impact.KnowledgeGain * multiplier))
	status.Collaboration = gamification.ClampAttribute(status.Collaboration + (impact.CollaborationImpact * multiplier))
	
	// Add experience
	expGained := impact.ExpGain * multiplier
	status.TotalExp += expGained
	
	// Calculate new level
	newLevel := gamification.CalculateLevel(status.TotalExp, s.config)
	if newLevel > status.Level {
		// Level up!
		status.Level = newLevel
		status.Experience = gamification.CalculateCurrentLevelExp(status.TotalExp, newLevel, s.config)
		
		// Could trigger level up notification here
	} else {
		status.Experience = gamification.CalculateCurrentLevelExp(status.TotalExp, status.Level, s.config)
	}
	
	// Update in database
	err = s.repo.Gamification().UpdateUserStatus(status)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	
	return nil
}

// ApplyAttributeModifier applies a temporary or permanent modifier
func (s *GamificationService) ApplyAttributeModifier(userID string, attribute string, value int, reason string, duration *time.Duration) error {
	modifier := &gamification.AttributeModifier{
		ID:        uuid.NewString(),
		UserID:    userID,
		Attribute: attribute,
		Value:     value,
		Reason:    reason,
	}
	
	if duration != nil {
		expiresAt := time.Now().Add(*duration)
		modifier.ExpiresAt = &expiresAt
	}
	
	err := s.repo.Gamification().CreateAttributeModifier(modifier)
	if err != nil {
		return fmt.Errorf("failed to create attribute modifier: %w", err)
	}
	
	return nil
}

// GetEffectiveStatus calculates status with all active modifiers applied
func (s *GamificationService) GetEffectiveStatus(userID string) (*gamification.UserStatus, error) {
	// Get base status
	status, err := s.GetOrCreateUserStatus(userID)
	if err != nil {
		return nil, err
	}
	
	// Get active modifiers
	modifiers, err := s.repo.Gamification().GetAttributeModifiers(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute modifiers: %w", err)
	}
	
	// Apply modifiers
	for _, mod := range modifiers {
		switch mod.Attribute {
		case "focus":
			status.Focus = gamification.ClampAttribute(status.Focus + mod.Value)
		case "productivity":
			status.Productivity = gamification.ClampAttribute(status.Productivity + mod.Value)
		case "creativity":
			status.Creativity = gamification.ClampAttribute(status.Creativity + mod.Value)
		case "stamina":
			status.Stamina = gamification.ClampAttribute(status.Stamina + mod.Value)
		case "knowledge":
			status.Knowledge = gamification.ClampAttribute(status.Knowledge + mod.Value)
		case "collaboration":
			status.Collaboration = gamification.ClampAttribute(status.Collaboration + mod.Value)
		}
	}
	
	return status, nil
}

// RestoreStamina gradually restores stamina over time
func (s *GamificationService) RestoreStamina(userID string, restoreAmount int) error {
	status, err := s.GetOrCreateUserStatus(userID)
	if err != nil {
		return err
	}
	
	status.Stamina = gamification.ClampAttribute(status.Stamina + restoreAmount)
	
	err = s.repo.Gamification().UpdateUserStatus(status)
	if err != nil {
		return fmt.Errorf("failed to update stamina: %w", err)
	}
	
	return nil
}