package backend

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/playwright-community/playwright-go"
)

func sleepCtx(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

func SearchForWebsites(ctx context.Context, city string, industry string, headless bool) ([]string, error) {
	if err := ctxErr(ctx); err != nil {
		return nil, err
	}

	websites := []string{}

	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("error running playwright: %v", err)
	}
	registerForceStop(func() { _ = pw.Stop() })
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		return nil, fmt.Errorf("error launching the browser: %v", err)
	}
	registerForceStop(func() { _ = browser.Close() })
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("error opening a page: %v", err)
	}
	registerForceStop(func() { _ = page.Close() })
	defer page.Close()

	if err := ctxErr(ctx); err != nil {
		return nil, err
	}

	_, err = page.Goto("https://www.gelbeseiten.de/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(20000),
	})
	if err != nil {
		return nil, fmt.Errorf("error opening the website: %v", err)
	}
	if err := ctxErr(ctx); err != nil {
		return nil, err
	}

	cookieLocator := page.Locator(".cmpboxbtn.cmpboxbtnyes.cmptxt_btn_yes")
	if err := cookieLocator.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}); err != nil {
		return nil, fmt.Errorf("error locating the accept cookie button: %v", err)
	}
	if err := ctxErr(ctx); err != nil {
		return nil, err
	}
	if err := cookieLocator.Click(); err != nil {
		return nil, fmt.Errorf("error accepting the cookies: %v", err)
	}
	if err := ctxErr(ctx); err != nil {
		return nil, err
	}

	industryLocator := page.Locator("#what_search")
	if err := industryLocator.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}); err != nil {
		return nil, fmt.Errorf("error locating the industry input: %v", err)
	}
	if err := industryLocator.Fill(industry); err != nil {
		return nil, fmt.Errorf("error filling in the industry: %v", err)
	}
	if err := ctxErr(ctx); err != nil {
		return nil, err
	}

	cityLocator := page.Locator("#where_search")
	if err := cityLocator.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}); err != nil {
		return nil, fmt.Errorf("error locating the city input: %v", err)
	}
	if err := cityLocator.Fill(city); err != nil {
		return nil, fmt.Errorf("error filling in the city: %v", err)
	}
	if err := ctxErr(ctx); err != nil {
		return nil, err
	}

	searchBtn := page.Locator(".gc-btn.gc-btn--black.gc-btn--l.search_go")
	if err := searchBtn.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}); err != nil {
		return nil, fmt.Errorf("error locating the search businesses button: %v", err)
	}
	if err := searchBtn.Click(); err != nil {
		return nil, fmt.Errorf("error clicking the search businesses button: %v", err)
	}
	if err := ctxErr(ctx); err != nil {
		return nil, err
	}

	rangeSlider := page.Locator("#suchradius_slider")
	if err := rangeSlider.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}); err != nil {
		return nil, fmt.Errorf("error locating the city range slider: %v", err)
	}
	if err := rangeSlider.Fill("50000"); err != nil {
		return nil, fmt.Errorf("error sliding the city range slider: %v", err)
	}
	if err := searchBtn.Click(); err != nil {
		return nil, fmt.Errorf("error clicking the search button after sliding: %v", err)
	}

	amountLocator := page.Locator("#mod-TrefferlisteInfo")
	if err := amountLocator.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}); err != nil {
		return nil, fmt.Errorf("error locating the amount of businesses text: %v", err)
	}

	amountStr, err := amountLocator.InnerText()
	if err != nil {
		return nil, fmt.Errorf("error reading amount text: %v", err)
	}

	amountInt, err := strconv.Atoi(amountStr)
	if err != nil {
		return nil, fmt.Errorf("error converting amount text to int: %v", err)
	}

	clicks := int(math.Ceil(float64(amountInt)/10)) - 5 // 50 already shown

	// was: page.WaitForTimeout(5000)
	if err := sleepCtx(ctx, 5*time.Second); err != nil {
		return nil, err
	}

	if amountInt > 50 {
		showMore := page.Locator("#mod-LoadMore--button")
		for i := 0; i < clicks; i++ {
			if err := ctxErr(ctx); err != nil {
				return nil, err
			}
			if err := showMore.WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}); err != nil {
				return nil, fmt.Errorf("error locating show more button: %v", err)
			}
			if err := showMore.Click(); err != nil {
				return nil, fmt.Errorf("error clicking show more button: %v", err)
			}
			// was: page.WaitForTimeout(1000)
			if err := sleepCtx(ctx, 1*time.Second); err != nil {
				return nil, err
			}
		}
	}

	websiteLinks := page.Locator(".mod-WebseiteKompakt__text")
	count, err := websiteLinks.Count()
	if err != nil {
		return nil, fmt.Errorf("error counting website links: %v", err)
	}

	for i := 0; i < count; i++ {
		if err := ctxErr(ctx); err != nil {
			// make "cancel" look clean
			if errors.Is(err, context.Canceled) {
				return websites, context.Canceled
			}
			return nil, err
		}

		link := websiteLinks.Nth(i)
		attr, err := link.GetAttribute("data-webseitelink")
		if err != nil {
			return nil, fmt.Errorf("error reading data-webseitelink for #%d: %v", i+1, err)
		}

		decoded, err := base64.StdEncoding.DecodeString(attr)
		if err != nil {
			return nil, fmt.Errorf("error decoding base64 url for #%d: %v", i+1, err)
		}

		websites = append(websites, string(decoded))
	}

	return websites, nil
}
