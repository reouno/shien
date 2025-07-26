package daemon

import (
	"context"
	"log"
	"time"
)

type Daemon struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func New() *Daemon {
	return &Daemon{}
}

func (d *Daemon) Start() error {
	d.ctx, d.cancel = context.WithCancel(context.Background())

	go d.run()

	return nil
}

func (d *Daemon) Stop() error {
	if d.cancel != nil {
		d.cancel()
	}
	return nil
}

func (d *Daemon) run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Println("Daemon is running...")

	for {
		select {
		case <-d.ctx.Done():
			log.Println("Daemon context cancelled")
			return
		case <-ticker.C:
			log.Println("Daemon heartbeat")
		}
	}
}