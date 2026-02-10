package backend

import (
	"strings"

	"github.com/playwright-community/playwright-go"
)

func TakeScreenshotOfWebsite(websiteUrl string, headless bool) (string, string) {
	pw, err := playwright.Run()
	if err != nil {
		return "", "error running playwright"
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		return "", "error launching the browser"
	}

	page, err := browser.NewPage()
	if err != nil {
		return "", "error opening a page"
	}

	_, err = page.Goto(websiteUrl, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	if err != nil {
		return "", websiteUrl
	}

	if strings.Contains(websiteUrl, "https") {
		_, err = page.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String("../website_screenshots/" + websiteUrl[8:] + ".png"),
		})
		if err != nil {
			return err.Error(), websiteUrl + "|screenshot not taken"
		}
	} else {
		_, err = page.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String("../website_screenshots/" + websiteUrl[7:] + ".png"),
		})
		if err != nil {
			return err.Error(), websiteUrl + "|screenshot not taken"
		}
	}

	return websiteUrl, ""
}
