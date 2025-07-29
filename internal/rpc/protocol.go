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