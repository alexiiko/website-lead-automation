<script>
  import { onMount } from "svelte"
  import Tabs from "./Tabs.svelte"

  import { SearchForWebsites } from "../wailsjs/go/main/App"
  import { TakeScreenshotOfWebsite } from "../wailsjs/go/main/App"
  import { ResetScreenshotsDir } from "../wailsjs/go/main/App"
  import { CancelCurrentJob } from "../wailsjs/go/main/App"

  import data from "./data.json"
  let availableIndustries = data[0].industries

  const LS = {
    lastIndustry: "lsa:lastSearchedIndustry",
    searchDone: "lsa:searchDone",
    searchAborted: "lsa:searchAborted",
    city: "lsa:city",
    industry: "lsa:selectedIndustry",
    headless: "lsa:headless"
  }

  let selectedIndustry = "Autoreparaturen"
  let lastSearchedIndustry = ""

  let statusText = ""
  let statusTextScreenshots = ""
  let statusTextUrls = ""
  let websiteUrls = []

  let city = "Merseburg"
  let headless = false

  let amountOfWebsites = 0
  let searchButtonActive = true
  let screenshotProgress = 0

  let searchDone = false
  let searchAborted = false

  function lsGet(key, fallback = "") {
    try {
      const v = localStorage.getItem(key)
      return v === null ? fallback : v
    } catch {
      return fallback
    }
  }

  function lsGetBool(key, fallback = false) {
    const v = lsGet(key, "")
    if (v === "") return fallback
    return v === "true"
  }

  function lsSet(key, value) {
    try {
      localStorage.setItem(key, String(value))
    } catch {
    }
  }

  function loadPersisted() {
    selectedIndustry = lsGet(LS.industry, selectedIndustry)
    city = lsGet(LS.city, city)
    headless = lsGetBool(LS.headless, headless)

    lastSearchedIndustry = lsGet(LS.lastIndustry, "")
    searchDone = lsGetBool(LS.searchDone, false)
    searchAborted = lsGetBool(LS.searchAborted, false)
  }

  function persistInputs() {
    lsSet(LS.industry, selectedIndustry)
    lsSet(LS.city, city)
    lsSet(LS.headless, headless)
  }

  function setLastSearchedIndustry(industry) {
    lsSet(LS.lastIndustry, industry)
  }

  function setSearchDone(v) {
    lsSet(LS.searchDone, v)
  }

  function setSearchAborted(v) {
    lsSet(LS.searchAborted, v)
  }

  onMount(() => {
    loadPersisted()
  })

  $: persistInputs()

  async function searchForWebsites() {
    if (city !== "") {
      try {
        websiteUrls = await SearchForWebsites(city, selectedIndustry, !headless)
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
      statusText = "Screenshots von Webseiten machen..."
    } catch (error) {
      statusTextScreenshots = ""
      statusText =
        `Fehler beim Erstellen eines Screenshots bei dieser Webseite:  ${url} \n` +
        error.toString()
    }
  }

  function abortSearch() {
    setSearchAborted(true)
    setSearchDone(false)

    statusText = ""
    statusTextUrls = ""
    statusTextScreenshots = ""
    websiteUrls = []
    amountOfWebsites = 0
    screenshotProgress = 0

    CancelCurrentJob()
  }

  async function main() {
    // new run
    setSearchAborted(false)
    setSearchDone(false)

    screenshotProgress = 0
    await ResetScreenshotsDir()

    setLastSearchedIndustry(selectedIndustry)

    statusText = ""
    statusTextScreenshots = ""
    statusTextUrls = ""
    websiteUrls = []
    amountOfWebsites = 0

    statusText = "URLs sammeln..."
    searchButtonActive = false
    await searchForWebsites()

    if (searchAborted) {
      statusText = ""
      searchButtonActive = true
      return
    }

    statusText = "Screenshots von Webseiten machen..."
    for (let index = 0; index < websiteUrls.length; index++) {
      if (searchAborted) break

      const websiteUrl = websiteUrls[index]
      statusTextScreenshots = `Bei Website: ${websiteUrl} \n${index + 1}/${websiteUrls.length}`
      await takeScreenShotsOfWebsite(websiteUrl)
      screenshotProgress++
    }

    if (!searchAborted) {
      statusTextScreenshots = `${screenshotProgress} Screenshots von ${websiteUrls.length} Webseiten gemacht!`
      setSearchDone(true)
    } else {
      statusTextScreenshots = "Suche abgebrochen."
      setSearchDone(false)
    }

    statusText = ""
    searchButtonActive = true
  }
</script>

<Tabs
  tabs={[
    { id: "website_search", label: "Webseiten suchen" },
    { id: "analyse_screenshots", label: "Bilder analysieren" }
  ]}
  canOpenAnalyse={searchDone}
>
  <div class="page" slot="website_search">
    <section class="card">
      <div class="form">
        <div class="grid">
          <div class="field">
            <label class="field-label" for="industries-select">🏭 Industrie:</label>
            <select
              bind:value={selectedIndustry}
              required
              disabled={!searchButtonActive}
              name="industries"
              id="industries-select"
            >
              {#each availableIndustries as industry}
                <option value={industry}>{industry}</option>
              {/each}
            </select>
            <label for="industries-select">⏳ Zuletzt gesucht: {lastSearchedIndustry}</label>
          </div>

          <div class="field">
            <label class="field-label" for="city-input">🏙️ Stadt:</label>
            <input
              id="city-input"
              disabled={!searchButtonActive}
              bind:value={city}
              placeholder="Berlin"
            />
          </div>
        </div>

        <div class="controls">
          <button class="btn btn--primary" on:click={main} disabled={!searchButtonActive}>
            <span class="btn-icon" aria-hidden="true">🔎</span>
            <span>Suchen</span>
          </button>

          <button
            class="btn btn--danger"
            disabled={searchButtonActive}
            on:click={abortSearch}
          >
            <span class="btn-icon" aria-hidden="true">❌</span>
            <span>Suche abbrechen</span>
          </button>

          <label class="toggle" title="Browser anzeigen (nicht headless)">
            <input
              class="toggle__input"
              type="checkbox"
              disabled={!searchButtonActive}
              bind:checked={headless}
            />
            <span class="toggle__track" aria-hidden="true"></span>
            <span class="toggle__label">🌐 Browser anzeigen</span>
          </label>
        </div>
      </div>

      <div id="status-texts">
        <pre>{statusTextUrls}</pre>
        <pre>{statusText}</pre>
        <pre>{statusTextScreenshots}</pre>
      </div>
    </section>
  </div>

  <div slot="analyse_screenshots">
    <h1>Analyse Screenshots</h1>
  </div>
</Tabs>

<style>
  :global(html, body) {
    height: 100%;
    margin: 0;
  }

  .page {
    min-height: 83vh;
    display: grid;
    place-items: center;
    box-sizing: border-box;
  }

  .card {
    width: min(760px, 100%);
    display: flex;
    flex-direction: column;
    gap: 18px;
    padding: 22px;
    border-radius: 14px;
    border: 1px solid rgba(0, 0, 0, 0.12);
    background: rgba(255, 255, 255, 0.9);
    box-shadow: 0 12px 30px rgba(0, 0, 0, 0.08);
  }

  .form {
    display: flex;
    flex-direction: column;
    gap: 14px;
  }

  .grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 14px;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 6px;
    align-items: flex-start;
  }

  .field-label {
    text-align: left;
    font-weight: bold;
    width: 100%;
    font-size: 13px;
    opacity: 0.85;
  }

  select,
  input {
    width: 100%;
    height: 40px;
    padding: 0 12px;
    border-radius: 10px;
    border: 1px solid rgba(0, 0, 0, 0.16);
    outline: none;
    background: white;
    box-sizing: border-box;
  }

  select:disabled,
  input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .controls {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: nowrap;
  }

  .btn {
    height: 42px;
    padding: 0 14px;
    border-radius: 10px;
    border: 1px solid rgba(0, 0, 0, 0.16);
    background: white;
    cursor: pointer;
    display: inline-flex;
    gap: 10px;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    white-space: nowrap;
  }

  .btn:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }

  .btn--primary {
    border-color: rgba(0, 0, 0, 0.18);
  }

  .toggle {
    margin-left: auto;
    display: inline-flex;
    align-items: center;
    gap: 10px;
    user-select: none;
    cursor: pointer;
    flex: 0 0 auto;
    white-space: nowrap;
  }

  .toggle__input {
    position: absolute;
    opacity: 0;
    pointer-events: none;
  }

  .toggle__track {
    width: 44px;
    height: 26px;
    border-radius: 999px;
    border: 1px solid rgba(0, 0, 0, 0.16);
    position: relative;
    background: rgba(214, 113, 147, 0.7);
    display: inline-block;
    flex: 0 0 auto;
  }

  .toggle__track::after {
    content: "";
    position: absolute;
    top: 50%;
    left: 4px;
    transform: translateY(-50%);
    width: 18px;
    height: 18px;
    border-radius: 999px;
    background: white;
    border: 1px solid rgba(0, 0, 0, 0.16);
    transition: left 150ms ease;
  }

  .toggle__input:checked + .toggle__track {
    background: rgba(46, 204, 113, 0.35);
    border-color: rgba(46, 204, 113, 0.7);
  }

  .toggle__input:checked + .toggle__track::after {
    left: 22px;
    border-color: rgba(46, 204, 113, 0.7);
  }

  .toggle__label {
    font-size: 13px;
    opacity: 0.9;
  }

  @media (max-width: 640px) {
    .grid {
      grid-template-columns: 1fr;
    }
    .controls {
      flex-wrap: wrap;
    }
    .toggle {
      margin-left: 0;
    }
  }
</style>
