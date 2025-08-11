package repository

import (
	"database/sql"
	"shien/internal/models/gamification"
	"time"
)

// GamificationRepo handles gamification-related database operations
type GamificationRepo struct {
	db *sql.DB
}

// NewGamificationRepo creates a new gamification repository
func NewGamificationRepo(db *sql.DB) *GamificationRepo {
	return &GamificationRepo{db: db}
}

// GetUserStatus retrieves the user's gamification status
func (r *GamificationRepo) GetUserStatus(userID string) (*gamification.UserStatus, error) {
	query := `
		SELECT user_id, level, experience, total_exp, 
		       focus, productivity, creativity, stamina, knowledge, collaboration,
		       updated_at, created_at
		FROM user_status
		WHERE user_id = ?
	`
	
	var status gamification.UserStatus
	err := r.db.QueryRow(query, userID).Scan(
		&status.UserID,
		&status.Level,
		&status.Experience,
		&status.TotalExp,
		&status.Focus,
		&status.Productivity,
		&status.Creativity,
		&status.Stamina,
		&status.Knowledge,
		&status.Collaboration,
		&status.UpdatedAt,
		&status.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &status, nil
}

// CreateUserStatus creates a new user status record
func (r *GamificationRepo) CreateUserStatus(status *gamification.UserStatus) error {
	query := `
		INSERT INTO user_status (
			user_id, level, experience, total_exp,
			focus, productivity, creativity, stamina, knowledge, collaboration,
			updated_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	status.CreatedAt = now
	status.UpdatedAt = now
	
	_, err := r.db.Exec(query,
		status.UserID,
		status.Level,
		status.Experience,
		status.TotalExp,
		status.Focus,
		status.Productivity,
		status.Creativity,
		status.Stamina,
		status.Knowledge,
		status.Collaboration,
		status.UpdatedAt,
		status.CreatedAt,
	)
	
	return err
}

// UpdateUserStatus updates an existing user status
func (r *GamificationRepo) UpdateUserStatus(status *gamification.UserStatus) error {
	query := `
		UPDATE user_status SET
			level = ?, experience = ?, total_exp = ?,
			focus = ?, productivity = ?, creativity = ?,
			stamina = ?, knowledge = ?, collaboration = ?,
			updated_at = ?
		WHERE user_id = ?
	`
	
	status.UpdatedAt = time.Now()
	
	_, err := r.db.Exec(query,
		status.Level,
		status.Experience,
		status.TotalExp,
		status.Focus,
		status.Productivity,
		status.Creativity,
		status.Stamina,
		status.Knowledge,
		status.Collaboration,
		status.UpdatedAt,
		status.UserID,
	)
	
	return err
}

// GetAttributeModifiers retrieves all active modifiers for a user
func (r *GamificationRepo) GetAttributeModifiers(userID string) ([]gamification.AttributeModifier, error) {
	query := `
		SELECT id, user_id, attribute, value, reason, expires_at, created_at
		FROM attribute_modifiers
		WHERE user_id = ?
		  AND (expires_at IS NULL OR expires_at > ?)
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var modifiers []gamification.AttributeModifier
	for rows.Next() {
		var mod gamification.AttributeModifier
		err := rows.Scan(
			&mod.ID,
			&mod.UserID,
			&mod.Attribute,
			&mod.Value,
			&mod.Reason,
			&mod.ExpiresAt,
			&mod.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		modifiers = append(modifiers, mod)
	}
	
	return modifiers, nil
}

// CreateAttributeModifier adds a new attribute modifier
func (r *GamificationRepo) CreateAttributeModifier(mod *gamification.AttributeModifier) error {
	query := `
		INSERT INTO attribute_modifiers (
			id, user_id, attribute, value, reason, expires_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	mod.CreatedAt = time.Now()
	
	_, err := r.db.Exec(query,
		mod.ID,
		mod.UserID,
		mod.Attribute,
		mod.Value,
		mod.Reason,
		mod.ExpiresAt,
		mod.CreatedAt,
	)
	
	return err
}

// CleanupExpiredModifiers removes expired modifiers
func (r *GamificationRepo) CleanupExpiredModifiers() error {
	query := `
		DELETE FROM attribute_modifiers
		WHERE expires_at IS NOT NULL AND expires_at <= ?
	`
	
	_, err := r.db.Exec(query, time.Now())
	return err
}