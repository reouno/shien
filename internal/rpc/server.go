package rpc

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
	
	"shien/internal/paths"
	"shien/internal/service"
)

// Server handles RPC requests
type Server struct {
	socketPath string
	listener   net.Listener
	services   *service.Services
	startedAt  time.Time
	mu         sync.RWMutex
	shutdown   chan struct{}
}

// NewServer creates a new RPC server
func NewServer(services *service.Services) (*Server, error) {
	socketPath := paths.SocketFile()
	
	// Remove existing socket file if it exists
	os.Remove(socketPath)
	
	return &Server{
		socketPath: socketPath,
		services:   services,
		startedAt:  time.Now(),
		shutdown:   make(chan struct{}),
	}, nil
}

// Start starts the RPC server
func (s *Server) Start() error {
	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("failed to listen on socket: %w", err)
	}
	
	// Set permissions so only the user can access
	if err := os.Chmod(s.socketPath, 0600); err != nil {
		listener.Close()
		return fmt.Errorf("failed to set socket permissions: %w", err)
	}
	
	s.listener = listener
	
	go s.acceptConnections()
	return nil
}

// Stop stops the RPC server
func (s *Server) Stop() error {
	close(s.shutdown)
	if s.listener != nil {
		s.listener.Close()
	}
	os.Remove(s.socketPath)
	return nil
}

func (s *Server) acceptConnections() {
	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.shutdown:
					return
				default:
					continue
				}
			}
			
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	
	var req Request
	if err := decoder.Decode(&req); err != nil {
		encoder.Encode(Response{
			Success: false,
			Error:   "invalid request format",
		})
		return
	}
	
	response := s.handleRequest(req)
	encoder.Encode(response)
}

func (s *Server) handleRequest(req Request) Response {
	switch req.Method {
	case MethodPing:
		return Response{
			Success: true,
			Data:    "pong",
		}
		
	case MethodGetStatus:
		return Response{
			Success: true,
			Data: Status{
				Running:   true,
				StartedAt: s.startedAt,
				Version:   "1.0.0",
			},
		}
		
	case MethodGetActivityLogs:
		var filter ActivityLogFilter
		if fromStr, ok := req.Params["from"].(string); ok {
			filter.From, _ = time.Parse(time.RFC3339, fromStr)
		}
		if toStr, ok := req.Params["to"].(string); ok {
			filter.To, _ = time.Parse(time.RFC3339, toStr)
		}
		
		logs, err := s.services.Activity.GetActivityLogs(filter.From, filter.To)
		if err != nil {
			return Response{
				Success: false,
				Error:   err.Error(),
			}
		}
		
		return Response{
			Success: true,
			Data:    logs,
		}
		
	case MethodGetConfig:
		return Response{
			Success: true,
			Data:    s.services.Config.GetConfig(),
		}
		
	case MethodGetGamificationStatus:
		userID := "default_user" // Default for now
		if id, ok := req.Params["user_id"].(string); ok {
			userID = id
		}
		
		status, err := s.services.Gamification.GetEffectiveStatus(userID)
		if err != nil {
			return Response{
				Success: false,
				Error:   err.Error(),
			}
		}
		
		// Calculate next level experience requirement
		config := s.services.Gamification.GetConfig()
		nextLevelExp := config.ExpForLevel(status.Level + 1)
		
		return Response{
			Success: true,
			Data: GamificationStatus{
				UserID:        status.UserID,
				Level:         status.Level,
				Experience:    status.Experience,
				TotalExp:      status.TotalExp,
				NextLevelExp:  nextLevelExp,
				Focus:         status.Focus,
				Productivity:  status.Productivity,
				Creativity:    status.Creativity,
				Stamina:       status.Stamina,
				Knowledge:     status.Knowledge,
				Collaboration: status.Collaboration,
				UpdatedAt:     status.UpdatedAt,
			},
		}
		
	case MethodGetGamificationDetails:
		userID := "default_user" // Default for now
		if id, ok := req.Params["user_id"].(string); ok {
			userID = id
		}
		
		status, err := s.services.Gamification.GetEffectiveStatus(userID)
		if err != nil {
			return Response{
				Success: false,
				Error:   err.Error(),
			}
		}
		
		// Get modifiers
		modifiers, err := s.services.Gamification.GetModifiers(userID)
		if err != nil {
			return Response{
				Success: false,
				Error:   err.Error(),
			}
		}
		
		// Convert modifiers to RPC format
		rpcModifiers := make([]AttributeModifier, len(modifiers))
		for i, mod := range modifiers {
			rpcModifiers[i] = AttributeModifier{
				Attribute: mod.Attribute,
				Value:     mod.Value,
				Reason:    mod.Reason,
				ExpiresAt: mod.ExpiresAt,
			}
		}
		
		// Get today's app usage
		today := time.Now()
		startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
		appUsage, _ := s.services.Activity.GetAppUsageSummary(startOfDay, today)
		
		// Calculate next level experience requirement
		config := s.services.Gamification.GetConfig()
		nextLevelExp := config.ExpForLevel(status.Level + 1)
		
		return Response{
			Success: true,
			Data: GamificationDetails{
				Status: &GamificationStatus{
					UserID:        status.UserID,
					Level:         status.Level,
					Experience:    status.Experience,
					TotalExp:      status.TotalExp,
					NextLevelExp:  nextLevelExp,
					Focus:         status.Focus,
					Productivity:  status.Productivity,
					Creativity:    status.Creativity,
					Stamina:       status.Stamina,
					Knowledge:     status.Knowledge,
					Collaboration: status.Collaboration,
					UpdatedAt:     status.UpdatedAt,
				},
				Modifiers:  rpcModifiers,
				RecentApps: appUsage,
			},
		}
		
	default:
		return Response{
			Success: false,
			Error:   fmt.Sprintf("unknown method: %s", req.Method),
		}
	}
}