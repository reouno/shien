package rpc

import (
	"time"
)

// Request represents an RPC request
type Request struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

// Response represents an RPC response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Methods
const (
	MethodPing            = "ping"
	MethodGetStatus       = "get_status"
	MethodGetActivityLogs = "get_activity_logs"
	MethodGetConfig       = "get_config"
	MethodUpdateConfig    = "update_config"
	MethodShutdown        = "shutdown"
	MethodGetGamificationStatus = "get_gamification_status"
	MethodGetGamificationDetails = "get_gamification_details"
)

// Status represents daemon status
type Status struct {
	Running   bool      `json:"running"`
	StartedAt time.Time `json:"started_at"`
	Version   string    `json:"version"`
}

// ActivityLogFilter for querying logs
type ActivityLogFilter struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// GamificationStatus represents user's gamification status
type GamificationStatus struct {
	UserID        string    `json:"user_id"`
	Level         int       `json:"level"`
	Experience    int       `json:"experience"`
	TotalExp      int       `json:"total_exp"`
	NextLevelExp  int       `json:"next_level_exp"`
	Focus         int       `json:"focus"`
	Productivity  int       `json:"productivity"`
	Creativity    int       `json:"creativity"`
	Stamina       int       `json:"stamina"`
	Knowledge     int       `json:"knowledge"`
	Collaboration int       `json:"collaboration"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GamificationDetails includes status with modifiers and recent activity
type GamificationDetails struct {
	Status       *GamificationStatus      `json:"status"`
	Modifiers    []AttributeModifier      `json:"modifiers,omitempty"`
	RecentApps   map[string]int          `json:"recent_apps,omitempty"` // App name -> minutes used today
}

// AttributeModifier for RPC
type AttributeModifier struct {
	Attribute string     `json:"attribute"`
	Value     int        `json:"value"`
	Reason    string     `json:"reason"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}