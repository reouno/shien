package daemon

import (
	"context"
	"time"
	
	"shien/internal/tray"
	"shien/internal/ui"
)

type Daemon struct {
	ctx     context.Context
	cancel  context.CancelFunc
	display *ui.Display
	tray    *tray.Tray
}

func New() *Daemon {
	return &Daemon{
		display: ui.NewDisplay(),
		tray:    tray.New(),
	}
}

func (d *Daemon) Start() error {
	d.ctx, d.cancel = context.WithCancel(context.Background())

	// Display startup message
	d.display.ShowBanner("Shien Daemon Started", "Supporting your knowledge work")
	d.display.ShowSuccess("Daemon is ready")
	
	// Send notification via system tray
	d.tray.SendNotification("Shien", "Support daemon started")

	// Start the daemon worker
	go d.run()

	return nil
}

func (d *Daemon) Stop() error {
	if d.cancel != nil {
		d.cancel()
	}
	// Stop the system tray
	d.tray.Stop()
	return nil
}

// StartTray starts the system tray UI (blocking call)
func (d *Daemon) StartTray() {
	d.tray.Start()
}

func (d *Daemon) run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	d.display.ShowInfo("Daemon monitoring started")

	for {
		select {
		case <-d.ctx.Done():
			d.display.ShowInfo("Daemon shutting down...")
			return
		case <-ticker.C:
			// Example: Show different messages based on time
			now := time.Now()
			if now.Minute()%5 == 0 {
				d.display.ShowAlert("Remember to take a break!")
				// Send notification via system tray
				d.tray.SendNotification("Break Reminder", "Time to rest your eyes and stretch!")
			} else {
				d.display.ShowInfo("System check - All systems operational")
			}
		}
	}
}