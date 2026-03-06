package backend

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func ResetScreenshotsDir() error {
	const screenshotsDir = "website_screenshots"
	return os.RemoveAll(screenshotsDir)
}

var invalidFileChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)

func safeFilenameFromURL(raw string) string {
	uStr := raw
	if !strings.HasPrefix(uStr, "http://") && !strings.HasPrefix(uStr, "https://") {
		uStr = "http://" + uStr
	}

	u, err := url.Parse(uStr)
	var host string
	if err != nil || u.Host == "" {
		host = raw
	} else {
		host = u.Hostname()
	}

	host = strings.TrimPrefix(host, "www.")
	name := "www." + host

	name = invalidFileChars.ReplaceAllString(name, "_")
	name = strings.Trim(name, "._ ")

	if name == "" || name == "www" {
		name = "www.website"
	}
	return name
}

func ctxErr(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return context.Canceled
		}
		return err
	}
	return nil
}

func TakeScreenshotOfWebsite(ctx context.Context, websiteURL string, headless bool) (string, error) {
	if err := ctxErr(ctx); err != nil {
		return "", err
	}

	pw, err := playwright.Run()
	if err != nil {
		return "", err
	}
	registerForceStop(func() { _ = pw.Stop() })
	defer pw.Stop()

	if err := ctxErr(ctx); err != nil {
		return "", err
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		return "", err
	}
	registerForceStop(func() { _ = browser.Close() })
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return "", err
	}
	registerForceStop(func() { _ = page.Close() })
	defer page.Close()

	if err := ctxErr(ctx); err != nil {
		return "", err
	}

	navURL := websiteURL
	if !strings.HasPrefix(navURL, "http://") && !strings.HasPrefix(navURL, "https://") {
		navURL = "http://" + navURL
	}

	_, err = page.Goto(navURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(20000),
	})
	if err != nil {
		return "", err
	}
	if err := ctxErr(ctx); err != nil {
		return "", err
	}

	outDir := "website_screenshots"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", err
	}

	file := safeFilenameFromURL(websiteURL) + ".png"
	outPath := filepath.Join(outDir, file)

	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(outPath),
		FullPage: playwright.Bool(false),
	})
	if err != nil {
		return "", err
	}
	if err := ctxErr(ctx); err != nil {
		return "", err
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(1 * time.Millisecond):
	}

	return outPath, nil
}
