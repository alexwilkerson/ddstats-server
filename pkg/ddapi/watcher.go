package ddapi

import (
	"time"
)

const (
	WatcherInterval = 10 * time.Second
)

// Watcher reaches out to the Devil Daggers backend and retrieves the top 100
// every WatcherInterval, then notifies the websocket hub
type Watcher struct {
	ticker *time.Ticker
	quit   chan struct{}
}

func NewWatcher() *Watcher {
	return &Watcher{
		quit: make(chan struct{}),
	}
}

func (w *Watcher) Start() {
	w.ticker = time.NewTicker(WatcherInterval)
	for {
		select {
		case <-w.ticker.C:
			// TODO: write the rest of the watcher script
		case <-w.quit:
			return
		}
	}
}

func (w *Watcher) Close() {
	w.ticker.Stop()
	close(w.quit)
}