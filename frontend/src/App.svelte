<script>
  import { SearchForWebsites } from "../wailsjs/go/main/App"
  import { TakeScreenshotOfWebsite } from "../wailsjs/go/main/App"

  let statusText = ""
  let statusTextScreenshots = ""
  let statusTextUrls = ""
  let websiteUrls = []

  let city = "Merseburg"
  let industry = "Gießereien"
  let headless = false

  let amountOfWebsites = 0

  let searchButtonActive = true

  async function searchForWebsites() {
    if (city != "" && industry != "") {
      try {

        websiteUrls = await SearchForWebsites(city, industry, !headless)
        amountOfWebsites = websiteUrls.length


        statusTextUrls = `${amountOfWebsites} Webseiten gefunden.`

      } catch (error) {
        statusTextUrls = ""
        statusText = "Fehler beim Suchen der Webseiten: \n" + error.toString()
      }
    } else {
      statusTextUrls = ""
      statusText = "Stadt und Industrie müssen ausgefüllt sein!"
    }
  }

  async function takeScreenShotsOfWebsite(url) {
    try {
      await TakeScreenshotOfWebsite(url, !headless)
    } catch (error) {
      statusTextScreenshots = ""
      statusText = `Fehler beim Erstellen eines Screenshots bei dieser Webseite:  ${url} \n` + error.toString()
    }
  }

  async function main() {
    statusText = "URLs sammeln..."
    searchButtonActive = false
    await searchForWebsites()

    statusText = "Screenshots von Webseiten machen..."
    for (let index = 0; index < websiteUrls.length; index++) {
      const websiteUrl = websiteUrls[index];
      statusTextScreenshots = `Bei Website: ${websiteUrl} \n` + `${index}/${websiteUrls.length}`
      await takeScreenShotsOfWebsite(websiteUrl)
    }
    statusTextScreenshots = ""
    statusText = "Screenshots von Webseiten gemacht!"

    searchButtonActive = true
  }
</script>

<main>
  <div id="inputs-container" class="inputs-container">
    <div id="inputs">
      <input id="industry-input" placeholder="Unternehmen" bind:value={industry}/>
      <input id="city-input" placeholder="Stadt" bind:value={city}/>

      <button on:click={main} disabled={!searchButtonActive}>Suchen</button>
      <button disabled={searchButtonActive}>Suche abbrechen</button>
    </div>

    <div id="headless-checkbox-container">
        <input type="checkbox" disabled={!searchButtonActive} bind:checked={headless}>
        <p>Browser anzeigen</p>
    </div> 
  </div>
  <p style="white-space: pre-line;">{statusTextUrls}</p>
  <p style="white-space: pre-line;">{statusTextScreenshots}</p>
  <p style="white-space: pre-line;">{statusText}</p>
</main>

<style>
  .inputs-container {
    display: flex;
  }
</style>
