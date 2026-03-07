package main

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func (a *App) GetScreenshotFilenames() ([]string, error) {
	dir := "website_screenshots"
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var filenames []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".png") {
			filenames = append(filenames, entry.Name())
		}
	}
	return filenames, nil
}

func (a *App) GetScreenshotBase64(filename string) (string, error) {
	path := filepath.Join("website_screenshots", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (a *App) WriteBusinessUrlToExcelDatabase(websiteUrl string) error {
	return backend.WriteBusinessUrlToExcelDatabase(websiteUrl)
}

const imageIndexFile = "image_index.txt"

func (a *App) SaveImageIndex(index int) error {
	return os.WriteFile(imageIndexFile, []byte(strconv.Itoa(index)), 0644)
}

func (a *App) LoadImageIndex() int {
	data, err := os.ReadFile(imageIndexFile)
	if err != nil {
		return 0
	}
	index, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}
	return index
}
