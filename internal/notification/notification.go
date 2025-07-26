package notification

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"
)

// Notifier interface defines the contract for notification implementations
type Notifier interface {
	Send(title, message string) error
	SendWithOptions(title, message string, opts Options) error
}

// Options for notifications
type Options struct {
	Sound    string
	Subtitle string
	Group    string
}

// Manager handles notifications with fallback strategies
type Manager struct {
	notifiers []Notifier
	mu        sync.RWMutex
}

// NewManager creates a new notification manager with default notifiers
func NewManager() *Manager {
	m := &Manager{
		notifiers: []Notifier{},
	}
	
	// Add notifiers in order of preference
	if runtime.GOOS == "darwin" {
		// 1. Try terminal-notifier first (best experience)
		if isTerminalNotifierAvailable() {
			m.notifiers = append(m.notifiers, &TerminalNotifier{})
		}
		// 2. Fallback to osascript (always available on macOS)
		m.notifiers = append(m.notifiers, &OSAScriptNotifier{})
	} else if runtime.GOOS == "linux" {
		// Linux: use notify-send
		m.notifiers = append(m.notifiers, &LinuxNotifier{})
	}
	
	return m
}

// Send sends a notification using the first available notifier
func (m *Manager) Send(title, message string) error {
	return m.SendWithOptions(title, message, Options{})
}

// SendWithOptions sends a notification with options
func (m *Manager) SendWithOptions(title, message string, opts Options) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if len(m.notifiers) == 0 {
		return fmt.Errorf("no notification backends available")
	}
	
	var lastErr error
	for _, notifier := range m.notifiers {
		if err := notifier.SendWithOptions(title, message, opts); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}
	
	return fmt.Errorf("all notifiers failed: %v", lastErr)
}

// TerminalNotifier uses terminal-notifier on macOS
type TerminalNotifier struct{}

func (t *TerminalNotifier) Send(title, message string) error {
	return t.SendWithOptions(title, message, Options{})
}

func (t *TerminalNotifier) SendWithOptions(title, message string, opts Options) error {
	args := []string{"-title", title, "-message", message}
	
	if opts.Sound != "" {
		args = append(args, "-sound", opts.Sound)
	}
	if opts.Subtitle != "" {
		args = append(args, "-subtitle", opts.Subtitle)
	}
	if opts.Group != "" {
		args = append(args, "-group", opts.Group)
	}
	
	return exec.Command("terminal-notifier", args...).Run()
}

// OSAScriptNotifier uses osascript (AppleScript) on macOS
type OSAScriptNotifier struct{}

func (o *OSAScriptNotifier) Send(title, message string) error {
	return o.SendWithOptions(title, message, Options{})
}

func (o *OSAScriptNotifier) SendWithOptions(title, message string, opts Options) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	
	if opts.Subtitle != "" {
		script = fmt.Sprintf(`display notification "%s" with title "%s" subtitle "%s"`, message, title, opts.Subtitle)
	}
	
	if opts.Sound != "" {
		script += fmt.Sprintf(` sound name "%s"`, opts.Sound)
	}
	
	return exec.Command("osascript", "-e", script).Run()
}

// LinuxNotifier uses notify-send on Linux
type LinuxNotifier struct{}

func (l *LinuxNotifier) Send(title, message string) error {
	return exec.Command("notify-send", title, message).Run()
}

func (l *LinuxNotifier) SendWithOptions(title, message string, opts Options) error {
	// notify-send doesn't support all options, but we can add urgency, etc. if needed
	return l.Send(title, message)
}

// Helper functions
func isTerminalNotifierAvailable() bool {
	return exec.Command("which", "terminal-notifier").Run() == nil
}