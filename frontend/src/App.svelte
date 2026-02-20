<script>
  import { SearchForWebsites } from "../wailsjs/go/main/App";

  let statusText = ""
  let statusTextScreenshots = ""
  let websiteUrls = []

  let city = "Berlin"
  let industry = "Elektriker"
  let headless = false

  let amountOfWebsites = 0

  async function searchForWebsites() {
    statusText = ""
    if (city != "" && industry != "") {
      try {
        websiteUrls = await SearchForWebsites(city, industry, !headless)
        amountOfWebsites = websiteUrls.length
        statusText = ""
      } catch (error) {
        statusText = "Fehler beim Suchen der Webseiten! \n" + "Fehlermeldung:  \n" + error.toString()
      }
    } else {
      statusText = "Stadt und Industrie müssen ausgefüllt sein!"
    }
  }

  function takeScreenShotsOfWebsites() {
    return null
  }

  function main() {
    let array = []
    statusText = "Fetching Urls..."
    searchForWebsites()
    statusText = "Taking screenshots of websites..."
    for (let index = 0; index < array.length; index++) {
      const element = array[index];
      takeScreenShotsOfWebsites()
    }
  }
</script>

<main>
  <div id="inputs-container" class="inputs-container">
    <div id="inputs">
      <input id="industry-input" placeholder="Unternehmen" bind:value={industry}/>
      <input id="city-input" placeholder="Stadt" bind:value={city}/>

      <button on:click={searchForWebsites}>Suchen</button>
    </div>

    <div id="headless-checkbox-container">
      <label id="headless-bool" class="switch" style="display: flex;">
        <input type="checkbox" bind:checked={headless}>
        <p>Browser anzeigen</p>
      </label>
    </div> 
    <p>{statusText}</p>
  </div>
</main>

<style>
  .inputs-container {
    display: flex;
  }
</style>
