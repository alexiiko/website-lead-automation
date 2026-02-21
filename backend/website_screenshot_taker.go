package backend

import (
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/playwright-community/playwright-go"
)


func ResetScreenshotsDir() error {
	const screenshotsDir = "website_screenshots"
	if err := os.RemoveAll(screenshotsDir); err != nil {
		return err
	}
	return nil
}

var invalidFileChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)

func safeFilenameFromURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		s := invalidFileChars.ReplaceAllString(raw, "_")
		return strings.Trim(s, "._ ")
	}

	host := u.Hostname()
	path := strings.Trim(u.EscapedPath(), "/")

	name := host
	if path != "" {
		name += "_" + strings.ReplaceAll(path, "/", "_")
	}

	name = invalidFileChars.ReplaceAllString(name, "_")
	name = strings.Trim(name, "._ ")

	if name == "" {
		name = "website"
	}
	return name
}

func TakeScreenshotOfWebsite(websiteURL string, headless bool) (string, error) {
	pw, err := playwright.Run()
	if err != nil {
		return "", err
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		return "", err
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return "", err
	}
	defer page.Close()

	_, err = page.Goto(websiteURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(20000),
	})
	if err != nil {
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
		FullPage: playwright.Bool(true),
	})
	if err != nil {
		return "", err
	}

	return outPath, nil
}
