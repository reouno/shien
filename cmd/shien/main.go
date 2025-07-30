package main

import (
	"flag"
	"log"

	"shien/internal/daemon"
	"shien/internal/paths"
)

func main() {
	// Parse command-line flags
	dataDir := flag.String("data-dir", "", "Custom data directory (default: ~/.config/shien)")
	flag.Parse()

	// Set custom data directory if provided
	if *dataDir != "" {
		if err := paths.SetDataDir(*dataDir); err != nil {
			log.Fatalf("Failed to set data directory: %v", err)
		}
		log.Printf("Using data directory: %s", paths.DataDir())
	}

	log.Println("Starting shien daemon...")

	d := daemon.New()
	
	// Start daemon in a goroutine
	go func() {
		if err := d.Start(); err != nil {
			log.Fatalf("Failed to start daemon: %v", err)
		}
	}()

	// Start system tray (this will block)
	d.StartTray()
}