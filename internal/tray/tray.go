package tray

import (
	"fmt"
	"github.com/getlantern/systray"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"shien/internal/notification"
	"shien/internal/version"
	"time"
)

type Tray struct {
	title         string
	tooltip       string
	notifications chan Notification
	quit          chan struct{}
	notifier      *notification.Manager
}

type Notification struct {
	Title   string
	Message string
	Time    time.Time
}

func New() *Tray {
	return &Tray{
		title:         "Shien",
		tooltip:       "Supporting knowledge workers",
		notifications: make(chan Notification, 100),
		quit:          make(chan struct{}),
		notifier:      notification.NewManager(),
	}
}

// Start initializes and runs the system tray
func (t *Tray) Start() {
	systray.Run(t.onReady, t.onExit)
}

// Stop gracefully shuts down the system tray
func (t *Tray) Stop() {
	close(t.quit)
	systray.Quit()
}

// SendNotification adds a notification to the queue and shows OS notification
func (t *Tray) SendNotification(title, message string) {
	// Add to internal queue
	select {
	case t.notifications <- Notification{
		Title:   title,
		Message: message,
		Time:    time.Now(),
	}:
	default:
		// Drop notification if channel is full
	}
	
	// Show OS notification using the notification manager
	t.notifier.SendWithOptions(title, message, notification.Options{
		Group: "shien-service",
	})
}

func (t *Tray) onReady() {
	// Set up the system tray icon and tooltip
	systray.SetTitle("支")  // Show "支" (support) character as icon
	systray.SetTooltip(t.tooltip)
	
	// Create menu items
	mStatus := systray.AddMenuItem("Status: Running", "Shien service status")
	mVersion := systray.AddMenuItem(fmt.Sprintf("Version: %s", version.GetVersion()), "Shien version")
	mVersion.Disable()
	systray.AddSeparator()
	
	// Recent activity menu
	mRecentActivity := systray.AddMenuItem("Recent Activity", "View recent activity logs")
	
	// Recent notifications submenu
	mNotifications := systray.AddMenuItem("Recent Notifications", "View recent notifications")
	mClearNotifications := systray.AddMenuItem("Clear Notifications", "Clear all notifications")
	
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit Shien")
	
	// Notification history
	var notificationHistory []Notification
	
	// Handle menu clicks and notifications
	go func() {
		for {
			select {
			case <-t.quit:
				return
				
			case notification := <-t.notifications:
				// Add to history
				notificationHistory = append(notificationHistory, notification)
				if len(notificationHistory) > 10 {
					notificationHistory = notificationHistory[1:]
				}
				
				// Update menu to show notification count
				count := len(notificationHistory)
				mNotifications.SetTitle(fmt.Sprintf("Recent Notifications (%d)", count))
				
				// Flash the menu item briefly to indicate new notification
				mNotifications.SetTitle(fmt.Sprintf("● Recent Notifications (%d)", count))
				time.Sleep(2 * time.Second)
				mNotifications.SetTitle(fmt.Sprintf("Recent Notifications (%d)", count))
				
			case <-mStatus.ClickedCh:
				// Toggle status display
				mStatus.SetTitle("Status: Running ✓")
				
			case <-mRecentActivity.ClickedCh:
				// Open terminal and run shienctl activity -today
				go func() {
					command := getShienCommand("activity -today")
					if err := openTerminalWithCommand(command); err != nil {
						t.SendNotification("Error", fmt.Sprintf("Failed to open terminal: %v", err))
					}
				}()
				
			case <-mNotifications.ClickedCh:
				// Show notification history (in real app, would open a window)
				if len(notificationHistory) == 0 {
					mNotifications.SetTitle("No recent notifications")
				} else {
					// For now, just update the title with the latest
					latest := notificationHistory[len(notificationHistory)-1]
					mNotifications.SetTitle(fmt.Sprintf("Latest: %s", latest.Title))
				}
				
			case <-mClearNotifications.ClickedCh:
				notificationHistory = nil
				mNotifications.SetTitle("Recent Notifications")
				
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func (t *Tray) onExit() {
	// Cleanup code here
}

// getShienCommand returns the appropriate shien command based on environment
func getShienCommand(args string) string {
	// In development, prefer local shien if it exists alongside shien-service
	if exePath, err := os.Executable(); err == nil {
		dir := filepath.Dir(exePath)
		shienPath := filepath.Join(dir, "shien")
		if _, err := os.Stat(shienPath); err == nil {
			// If running from source directory (not installed path)
			if !strings.Contains(dir, "/usr/") && !strings.Contains(dir, "/opt/") && !strings.Contains(dir, "go/bin") {
				return fmt.Sprintf("%s %s", shienPath, args)
			}
		}
	}
	
	// Otherwise use shien from PATH
	return fmt.Sprintf("shien %s", args)
}

// openTerminalWithCommand opens a terminal and runs the specified command
func openTerminalWithCommand(command string) error {
	switch runtime.GOOS {
	case "darwin":
		// macOS - Use osascript to open Terminal app
		script := fmt.Sprintf(`tell application "Terminal"
			do script "%s"
			activate
		end tell`, command)
		cmd := exec.Command("osascript", "-e", script)
		return cmd.Start()
	case "linux":
		// Try common terminal emulators
		terminals := []string{"gnome-terminal", "konsole", "xterm", "xfce4-terminal"}
		for _, term := range terminals {
			if _, err := exec.LookPath(term); err == nil {
				var cmd *exec.Cmd
				switch term {
				case "gnome-terminal":
					cmd = exec.Command(term, "--", "bash", "-c", command+"; read -p 'Press Enter to close...'")
				case "konsole":
					cmd = exec.Command(term, "-e", "bash", "-c", command+"; read -p 'Press Enter to close...'")
				default:
					cmd = exec.Command(term, "-e", "bash", "-c", command+"; read -p 'Press Enter to close...'")
				}
				return cmd.Start()
			}
		}
		return fmt.Errorf("no terminal emulator found")
	case "windows":
		// Windows - Use cmd.exe
		cmd := exec.Command("cmd", "/c", "start", "cmd", "/k", command)
		return cmd.Start()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}