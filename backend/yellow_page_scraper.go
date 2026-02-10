package backend

import (
	"encoding/base64"
	"fmt"
	"math"
	"strconv"

	"github.com/playwright-community/playwright-go"
)

func SearchForWebsites(city string, industry string, headless bool) ([]string, error) {
	websites := []string{};

	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("error running playwright: %v", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		return nil, fmt.Errorf("error launching the browser: %v", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("error opening a page: %v", err)
	}

	// open yellow pages website 
	_, err = page.Goto("https://www.gelbeseiten.de/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
  if err != nil {
		return nil, fmt.Errorf("error opening the website: %v", err)
	}

	// accept cookies
	cookieLocator := page.Locator(".cmpboxbtn.cmpboxbtnyes.cmptxt_btn_yes")
	err = cookieLocator.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	})
	if err != nil {
		return nil, fmt.Errorf("error locating the accept cookie button: %v", err)
	}

  err = cookieLocator.Click()
	if err != nil {
		return nil, fmt.Errorf("error accepting the cookies: %v", err)
	}

	// fill out industry
  industryLocator := page.Locator("#what_search")
	err = industryLocator.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	})
	if err != nil {
		return nil, fmt.Errorf("error locating the industry input: %v", err)
	}

	err = industryLocator.Fill(industry)
	if err != nil {
		return nil, fmt.Errorf("error filling in the industry in the industry input: %v", err)
	}

	// fill out city
	cityLocator := page.Locator("#where_search")
	err = cityLocator.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	})
	if err != nil {
		return nil, fmt.Errorf("error locating the city input: %v", err)
	}

	err = cityLocator.Fill(city)
	if err != nil {
		return nil, fmt.Errorf("error filling in the city in the city input: %v", err)
	}

  // click the search businesses button
	searchBusinessesButtonLocator := page.Locator(".gc-btn.gc-btn--black.gc-btn--l.search_go")
	err = searchBusinessesButtonLocator.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	})
	if err != nil {
		return nil, fmt.Errorf("error locating the search businesses button: %v", err)
	}

	err = searchBusinessesButtonLocator.Click()
	if err != nil {
		return nil, fmt.Errorf("error clicking the search businesses button: %v", err)
	}

	// maximize the range for the city
	rangeSliderLocator := page.Locator("#suchradius_slider")
	err = rangeSliderLocator.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	})
  if err != nil {
		return nil, fmt.Errorf("error locating the city range slider: %v", err)
	}

	err = rangeSliderLocator.Fill("50000")
	if err != nil {
		return nil, fmt.Errorf("error sliding the city range slider: %v", err)
	}
	  // click the search businesses button again as the page does not update itself
		err = searchBusinessesButtonLocator.Click()
		if err != nil {
      return nil, fmt.Errorf("error clicking the search businesses button after sliding the city range slider: %v", err)
		}

	// click the show more businesses button amountOfBusinesses/10 times
	amountOfBusinessesLocator := page.Locator("#mod-TrefferlisteInfo")
  err = amountOfBusinessesLocator.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	})
	if err != nil {
		return nil, fmt.Errorf("error locating the amount of businesses text: %v", err)
	}

	amountOfBusinessesString, err := amountOfBusinessesLocator.InnerText()
	if err != nil {
		return nil, fmt.Errorf("error retrieving the text of the amount of business locator: %v", err)
	}

	amountOfShowMoreBusinessesInteger, err := strconv.Atoi(amountOfBusinessesString)
	if err != nil {
		return nil, fmt.Errorf("error converting the amountOfBusinessesString to an integer: %v", err)
	}

	amountOfShowMoreBusinessesClicks := int(math.Ceil(float64(amountOfShowMoreBusinessesInteger) / 10)) - 5 // 50 businesses are already shown

	page.WaitForTimeout(5000)

	showMoreBusinessesButtonLocator := page.Locator("#mod-LoadMore--button")
	for i := 0; i < int(amountOfShowMoreBusinessesClicks); i++ {
		err = showMoreBusinessesButtonLocator.WaitFor(playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateVisible,
		})
		if err != nil {
			return nil, fmt.Errorf("error locating the show more businesses button: %v", err)
		}

    err = showMoreBusinessesButtonLocator.Click()
		if err != nil {
			return nil, fmt.Errorf("error clicking the show more businesses button: %v", err)
		}
		page.WaitForTimeout(1000)
	}

	// get website links
	websiteLinksLocator := page.Locator(".mod-WebseiteKompakt__text")
	amountOfWebsiteLinks, err := websiteLinksLocator.Count()
	if err != nil {
		return nil, fmt.Errorf("error counting the links to the websites of the businesses: %v", err)
	}

  for i := 0; i < amountOfWebsiteLinks; i++ {
		businessWebsiteDataLink := websiteLinksLocator.Nth(i)

		websiteUrlLocator, err := businessWebsiteDataLink.GetAttribute("data-webseitelink")
		if err != nil {
			return nil, fmt.Errorf("error retrieving the website link of the %d business: %v", i + 1, err)
		}

		businessWebsiteUrl, err := base64.StdEncoding.DecodeString(websiteUrlLocator)
		if err != nil {
			return nil, fmt.Errorf("error encoding the base64 website of the %d business to an url: %v", i+1, err)
		}
		websites = append(websites, string(businessWebsiteUrl))
	}

	return websites,nil
}
