<script>
  import { SearchForWebsites } from "../wailsjs/go/main/app";

  let statusText = ""
  let statusTextScreenshots = ""
  let websiteUrls = null

  let city = ""
  let industry = ""
  let headless = false

  let amountOfWebsites = 0

  function searchForWebsites() {
    if (city != "" && industry != "") {
      SearchForWebsites(city, industry, !headless).then((result) => { // ! because when the value is true the browser gets launched in headless mode when it should be launched in head mode
        websiteUrls = result
        amountOfWebsites = websiteUrls.length
      })
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
        <input type="checkbox" bind:value={headless}>
        <p>Browser anzeigen</p>
      </label>
    </div> 
  </div>
</main>

<style>
  .inputs-container {
    display: flex;
  }
</style>
