package main

import (
	"context"
	"website-lead-automation-go/backend"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// send the website urls to the frontend
func (a *App) SearchForWebsites(city string, industry string, headless bool) ([]string, error) {
	return backend.SearchForWebsites(city, industry, headless)
}

func (a *App) TakeScreenshotOfWebsite(url string, headless bool) (string, error) {
	return backend.TakeScreenshotOfWebsite(url, headless)
}
