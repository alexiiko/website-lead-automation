package backend

import (
	"strings"

	"github.com/playwright-community/playwright-go"
)

func TakeScreenshotOfWebsite(websiteUrl string, headless bool) (string, error) {
	pw, err := playwright.Run()
	if err != nil {
		return "", err
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		return "", err
	}

	page, err := browser.NewPage()
	if err != nil {
		return "", err
	}

	_, err = page.Goto(websiteUrl, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	if err != nil {
		return "", err
	}

	if strings.Contains(websiteUrl, "https") {
		_, err = page.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String("../website_screenshots/" + websiteUrl[8:] + ".png"),
		})
		if err != nil {
			return websiteUrl, err
		}
	} else {
		_, err = page.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String("../website_screenshots/" + websiteUrl[7:] + ".png"),
		})
		if err != nil {
			return websiteUrl, err
		}
	}

	page.Close()
	return websiteUrl, nil
}
