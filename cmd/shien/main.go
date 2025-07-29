package main

import (
	"log"

	"shien/internal/daemon"
)

func main() {
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