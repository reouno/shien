package rpc

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

// Client connects to the RPC server
type Client struct {
	socketPath string
}

// NewClient creates a new RPC client
func NewClient() (*Client, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	socketPath := filepath.Join(homeDir, ".config", "shien", "shien.sock")
	
	return &Client{
		socketPath: socketPath,
	}, nil
}

// Call makes an RPC call to the server
func (c *Client) Call(method string, params map[string]interface{}) (*Response, error) {
	// Check if socket exists
	if _, err := os.Stat(c.socketPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("daemon not running (socket not found)")
	}
	
	// Connect to socket
	conn, err := net.DialTimeout("unix", c.socketPath, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to daemon: %w", err)
	}
	defer conn.Close()
	
	// Set timeout for the entire operation
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	
	// Send request
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(Request{
		Method: method,
		Params: params,
	}); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	
	// Read response
	decoder := json.NewDecoder(conn)
	var response Response
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	return &response, nil
}

// Ping checks if the daemon is running
func (c *Client) Ping() error {
	resp, err := c.Call(MethodPing, nil)
	if err != nil {
		return err
	}
	
	if !resp.Success {
		return fmt.Errorf("ping failed: %s", resp.Error)
	}
	
	return nil
}

// GetStatus gets the daemon status
func (c *Client) GetStatus() (*Status, error) {
	resp, err := c.Call(MethodGetStatus, nil)
	if err != nil {
		return nil, err
	}
	
	if !resp.Success {
		return nil, fmt.Errorf("failed to get status: %s", resp.Error)
	}
	
	// Convert response data to Status
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}
	
	var status Status
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, err
	}
	
	return &status, nil
}