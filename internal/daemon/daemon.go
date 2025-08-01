package daemon

import (
	"context"
	"log"
	"time"

	"shien/internal/config"
	"shien/internal/database"
	"shien/internal/rpc"
	"shien/internal/service"
	"shien/internal/tray"
	"shien/internal/ui"
)

type Daemon struct {
	ctx       context.Context
	cancel    context.CancelFunc
	display   *ui.Display
	tray      *tray.Tray
	config    *config.Manager
	db        *database.DB
	repo      *database.Repository
	services  *service.Services
	rpcServer *rpc.Server
}

func New() *Daemon {
	// Initialize config manager
	configMgr, err := config.NewManager()
	if err != nil {
		log.Printf("Failed to initialize config: %v, using defaults", err)
		configMgr = &config.Manager{}
	}

	// Initialize database
	db, err := database.New()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create repository
	repo := database.NewRepository(db)
	
	// Create service layer
	services := service.NewServices(repo, configMgr)
	
	// Create RPC server
	rpcServer, err := rpc.NewServer(services)
	if err != nil {
		log.Printf("Failed to create RPC server: %v", err)
	}

	return &Daemon{
		display:   ui.NewDisplay(),
		tray:      tray.New(),
		config:    configMgr,
		db:        db,
		repo:      repo,
		services:  services,
		rpcServer: rpcServer,
	}
}

func (d *Daemon) Start() error {
	d.ctx, d.cancel = context.WithCancel(context.Background())

	// Display startup message
	d.display.ShowBanner("Shien Service Started", "Supporting your knowledge work")
	d.display.ShowSuccess("Daemon is ready")

	// Show config and database locations
	if d.config != nil {
		d.display.ShowInfo("Config: " + d.config.ConfigPath())
	}
	if d.db != nil {
		d.display.ShowInfo("Database: " + d.db.Path())
	}

	// Don't record initial activity - wait for the next 5-minute interval

	// Start RPC server
	if d.rpcServer != nil {
		if err := d.rpcServer.Start(); err != nil {
			log.Printf("Failed to start RPC server: %v", err)
		} else {
			d.display.ShowInfo("RPC server listening on Unix socket")
		}
	}

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

	// Stop RPC server
	if d.rpcServer != nil {
		d.rpcServer.Stop()
	}

	// Close database
	if d.db != nil {
		d.db.Close()
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
	d.display.ShowInfo("Daemon monitoring started")

	// Calculate time until next 5-minute interval
	now := time.Now()
	// Round up to next 5-minute mark
	nextInterval := now.Truncate(5 * time.Minute)
	if nextInterval.Before(now) || nextInterval.Equal(now) {
		nextInterval = nextInterval.Add(5 * time.Minute)
	}
	waitDuration := nextInterval.Sub(now)
	
	d.display.ShowInfo("Waiting until next 5-minute interval: " + nextInterval.Format("15:04:05"))
	
	// Wait until the next 5-minute interval
	timer := time.NewTimer(waitDuration)
	defer timer.Stop()
	
	select {
	case <-d.ctx.Done():
		d.display.ShowInfo("Daemon shutting down...")
		return
	case <-timer.C:
		// Record first activity at the aligned time
		if d.services != nil {
			if err := d.services.Activity.RecordActivity(); err != nil {
				log.Printf("Failed to record activity: %v", err)
			} else {
				d.display.ShowInfo("Activity recorded at " + time.Now().Format("15:04:05"))
			}
		}
	}
	
	// Now start regular 5-minute ticker
	activityTicker := time.NewTicker(5 * time.Minute)
	defer activityTicker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			d.display.ShowInfo("Daemon shutting down...")
			return
		case <-activityTicker.C:
			// Record activity
			if d.services != nil {
				if err := d.services.Activity.RecordActivity(); err != nil {
					log.Printf("Failed to record activity: %v", err)
				} else {
					d.display.ShowInfo("Activity recorded at " + time.Now().Format("15:04:05"))
				}
			}
		}
	}
}
