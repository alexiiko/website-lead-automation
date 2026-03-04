package main

import (
	"context"
	"sync"

	"website-lead-automation-go/backend"
)

type App struct {
	ctx context.Context

	mu          sync.Mutex
	currentCtx  context.Context
	currentStop context.CancelFunc
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) { a.ctx = ctx }

// Start (or restart) a cancelable job context
func (a *App) beginJob() context.Context {
	a.mu.Lock()
	defer a.mu.Unlock()

	// cancel any previous job
	if a.currentStop != nil {
		a.currentStop()
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.currentCtx = ctx
	a.currentStop = cancel
	return ctx
}

// Called from frontend when user clicks "Abbrechen"
func (a *App) CancelCurrentJob() {
	a.mu.Lock()
	cancel := a.currentStop
	a.mu.Unlock()
	if cancel != nil {
		cancel()
	}

	// Optional: also force-close Playwright resources ASAP
	backend.ForceStopPlaywright()
}

func (a *App) SearchForWebsites(city, industry string, headless bool) ([]string, error) {
	ctx := a.beginJob()
	return backend.SearchForWebsites(ctx, city, industry, headless)
}

func (a *App) TakeScreenshotOfWebsite(url string, headless bool) (string, error) {
	// If screenshots are part of same “job”, don’t call beginJob() here.
	// Use the current job ctx instead:
	a.mu.Lock()
	ctx := a.currentCtx
	a.mu.Unlock()

	// Fallback if someone calls screenshot directly without starting a job:
	if ctx == nil {
		ctx = a.beginJob()
	}

	return backend.TakeScreenshotOfWebsite(ctx, url, headless)
}

func (a *App) ResetScreenshotsDir() error {
	return backend.ResetScreenshotsDir()
}
