package tray

import (
	"fmt"
	"github.com/getlantern/systray"
	"shien/internal/notification"
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
		Group: "shien",
	})
}

func (t *Tray) onReady() {
	// Set up the system tray icon and tooltip
	systray.SetTitle("支")  // Show "支" (support) character as icon
	systray.SetTooltip(t.tooltip)
	
	// Create menu items
	mStatus := systray.AddMenuItem("Status: Running", "Shien daemon status")
	systray.AddSeparator()
	
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