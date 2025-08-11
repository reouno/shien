package gamification

import "time"

// UserStatus represents the gamification status of a user
type UserStatus struct {
	UserID      string    `json:"user_id" db:"user_id"`
	Level       int       `json:"level" db:"level"`
	Experience  int       `json:"experience" db:"experience"`
	TotalExp    int       `json:"total_exp" db:"total_exp"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	
	// Work-related attributes (no upper limit)
	Focus       int       `json:"focus" db:"focus"`         // 集中力
	Productivity int     `json:"productivity" db:"productivity"` // 生産性
	Creativity  int       `json:"creativity" db:"creativity"`   // 創造性
	Stamina     int       `json:"stamina" db:"stamina"`       // スタミナ
	Knowledge   int       `json:"knowledge" db:"knowledge"`     // 知識
	Collaboration int    `json:"collaboration" db:"collaboration"` // 協調性
}

// StatusConfig defines the rules for status calculations
type StatusConfig struct {
	// Experience points needed for each level
	// Level 1: 0-99, Level 2: 100-299, Level 3: 300-599, etc.
	ExpForLevel func(level int) int
}

// DefaultStatusConfig returns the default configuration
func DefaultStatusConfig() *StatusConfig {
	return &StatusConfig{
		ExpForLevel: func(level int) int {
			// Exponential growth: 100, 200, 300, 500, 800, ...
			if level <= 1 {
				return 0
			}
			base := 100
			for i := 2; i < level; i++ {
				base = int(float64(base) * 1.5)
			}
			return base
		},
	}
}

// AttributeModifier represents a temporary or permanent modifier to an attribute
type AttributeModifier struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Attribute   string    `json:"attribute" db:"attribute"`     // "focus", "productivity", etc.
	Value       int       `json:"value" db:"value"`             // Can be positive or negative
	Reason      string    `json:"reason" db:"reason"`           // Why this modifier was applied
	ExpiresAt   *time.Time `json:"expires_at" db:"expires_at"`  // NULL for permanent
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ActivityImpact defines how different activities affect status attributes
type ActivityImpact struct {
	AppName      string `json:"app_name"`
	Category     string `json:"category"`      // "development", "communication", "learning", etc.
	FocusImpact  int    `json:"focus_impact"`  // -10 to +10
	ProductivityImpact int `json:"productivity_impact"`
	CreativityImpact   int `json:"creativity_impact"`
	StaminaCost  int    `json:"stamina_cost"` // 0 to 10 (always positive, represents drain)
	KnowledgeGain int   `json:"knowledge_gain"` // 0 to 5
	CollaborationImpact int `json:"collaboration_impact"`
	ExpGain      int    `json:"exp_gain"`      // Experience points gained
}

// PredefinedActivityImpacts returns common activity impacts
func PredefinedActivityImpacts() map[string]ActivityImpact {
	return map[string]ActivityImpact{
		"Code Editor": {
			AppName:      "Code Editor",
			Category:     "development",
			FocusImpact:  8,
			ProductivityImpact: 10,
			CreativityImpact: 6,
			StaminaCost:  3,
			KnowledgeGain: 3,
			CollaborationImpact: 0,
			ExpGain:      15,
		},
		"Terminal": {
			AppName:      "Terminal",
			Category:     "development",
			FocusImpact:  6,
			ProductivityImpact: 8,
			CreativityImpact: 3,
			StaminaCost:  2,
			KnowledgeGain: 2,
			CollaborationImpact: 0,
			ExpGain:      10,
		},
		"Slack": {
			AppName:      "Slack",
			Category:     "communication",
			FocusImpact:  -3,
			ProductivityImpact: 3,
			CreativityImpact: 2,
			StaminaCost:  1,
			KnowledgeGain: 1,
			CollaborationImpact: 8,
			ExpGain:      5,
		},
		"Browser": {
			AppName:      "Browser",
			Category:     "mixed",
			FocusImpact:  0,
			ProductivityImpact: 5,
			CreativityImpact: 4,
			StaminaCost:  2,
			KnowledgeGain: 4,
			CollaborationImpact: 2,
			ExpGain:      8,
		},
		"Documentation": {
			AppName:      "Documentation",
			Category:     "learning",
			FocusImpact:  5,
			ProductivityImpact: 6,
			CreativityImpact: 2,
			StaminaCost:  2,
			KnowledgeGain: 8,
			CollaborationImpact: 0,
			ExpGain:      10,
		},
		"Design Tool": {
			AppName:      "Design Tool",
			Category:     "creative",
			FocusImpact:  7,
			ProductivityImpact: 7,
			CreativityImpact: 10,
			StaminaCost:  3,
			KnowledgeGain: 2,
			CollaborationImpact: 2,
			ExpGain:      12,
		},
		"Email": {
			AppName:      "Email",
			Category:     "communication",
			FocusImpact:  -2,
			ProductivityImpact: 4,
			CreativityImpact: 1,
			StaminaCost:  2,
			KnowledgeGain: 1,
			CollaborationImpact: 6,
			ExpGain:      5,
		},
		"Video Conference": {
			AppName:      "Video Conference",
			Category:     "communication",
			FocusImpact:  -5,
			ProductivityImpact: 3,
			CreativityImpact: 3,
			StaminaCost:  4,
			KnowledgeGain: 2,
			CollaborationImpact: 10,
			ExpGain:      8,
		},
	}
}

// CalculateLevel calculates the level based on total experience (no upper limit)
func CalculateLevel(totalExp int, config *StatusConfig) int {
	level := 1
	for {
		requiredExp := config.ExpForLevel(level + 1)
		if totalExp < requiredExp {
			break
		}
		level++
	}
	return level
}

// CalculateCurrentLevelExp calculates experience within the current level
func CalculateCurrentLevelExp(totalExp int, level int, config *StatusConfig) int {
	if level <= 1 {
		return totalExp
	}
	levelStartExp := config.ExpForLevel(level)
	return totalExp - levelStartExp
}

// ClampAttribute ensures an attribute value doesn't go negative (no upper limit)
func ClampAttribute(value int) int {
	if value < 0 {
		return 0
	}
	return value
}