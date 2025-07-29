package rpc

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
	
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	socketDir := filepath.Join(homeDir, ".config", "shien")
	socketPath := filepath.Join(socketDir, "shien.sock")
	
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
		
	default:
		return Response{
			Success: false,
			Error:   fmt.Sprintf("unknown method: %s", req.Method),
		}
	}
}