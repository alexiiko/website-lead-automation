package backend

import "sync"

var (
	stopMu     sync.Mutex
	forceStops []func()
)

// Register a cleanup that should run on cancel (close browser, stop pw, etc.)
func registerForceStop(fn func()) {
	stopMu.Lock()
	defer stopMu.Unlock()
	forceStops = append(forceStops, fn)
}

// Called by App.CancelCurrentJob()
func ForceStopPlaywright() {
	stopMu.Lock()
	fns := forceStops
	forceStops = nil
	stopMu.Unlock()

	for _, fn := range fns {
		// best-effort
		func() { defer func() { _ = recover() }(); fn() }()
	}
}
