package notify

import (
	"shien/internal/notification"
	"sync"
)

var (
	globalManager *notification.Manager
	once          sync.Once
)

// Send sends a simple notification using the global manager
func Send(title, message string) error {
	return Manager().Send(title, message)
}

// SendWithSound sends a notification with a sound
func SendWithSound(title, message, sound string) error {
	return Manager().SendWithOptions(title, message, notification.Options{
		Sound: sound,
	})
}

// Manager returns the global notification manager
func Manager() *notification.Manager {
	once.Do(func() {
		globalManager = notification.NewManager()
	})
	return globalManager
}