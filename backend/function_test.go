package backend 

import "testing"

func TestScrape(t *testing.T) {
	result, err := SearchForWebsites("Berlin", "Autoreparaturen", true)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}

func TestWebsitescreenshot(t *testing.T) {
	result, err := TakeScreenshotOfWebsite("https://www.gaebel-berlin.de", true)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}
