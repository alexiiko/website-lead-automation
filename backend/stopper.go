package backend

import "sync"

var (
	stopMu     sync.Mutex
	forceStops []func()
)

func registerForceStop(fn func()) {
	stopMu.Lock()
	defer stopMu.Unlock()
	forceStops = append(forceStops, fn)
}

func ForceStopPlaywright() {
	stopMu.Lock()
	fns := forceStops
	forceStops = nil
	stopMu.Unlock()

	for _, fn := range fns {
		func() { defer func() { _ = recover() }(); fn() }()
	}
}
