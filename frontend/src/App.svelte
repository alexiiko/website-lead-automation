<script>
  import { onMount } from "svelte";
  import Tabs from "./Tabs.svelte";

  import { SearchForWebsites } from "../wailsjs/go/main/App";
  import { TakeScreenshotOfWebsite } from "../wailsjs/go/main/App";
  import { ResetScreenshotsDir } from "../wailsjs/go/main/App";
  import { CancelCurrentJob } from "../wailsjs/go/main/App";
  import { GetScreenshotFilenames } from "../wailsjs/go/main/App";
  import { GetScreenshotBase64 } from "../wailsjs/go/main/App";
  import { WriteBusinessUrlToExcelDatabase } from "../wailsjs/go/main/App";
  import { SaveImageIndex } from "../wailsjs/go/main/App";
  import { LoadImageIndex } from "../wailsjs/go/main/App";

  import data from "./data.json";
  let availableIndustries = data[0].industries;
  let availableCities = data[1].cities;

  const LS = {
    lastIndustry: "lsa:lastSearchedIndustry",
    lastCity: "lsa:lastSearchedCity",
    searchDone: "lsa:searchDone",
    searchAborted: "lsa:searchAborted",
    city: "lsa:city",
    industry: "lsa:selectedIndustry",
    headless: "lsa:headless",
  };

  let selectedIndustry = "Autoreparaturen";
  let lastSearchedIndustry = "";

  let city = "Flensburg";
  let lastSearchedCity = "";

  let statusText = "";
  let statusTextScreenshots = "";
  let statusTextUrls = "";
  let websiteUrls = [];

  let headless = false;

  let amountOfWebsites = 0;
  let searchButtonActive = true;
  let screenshotProgress = 0;

  let searchDone = false;
  let searchAborted = false;

  let currentTab = "website_search";

  let screenshotFiles = [];
  let currentScreenshotIndex = 0;
  let currentScreenshotBase64 = "";
  let amountSavedToExcel = 0;

  function lsGet(key, fallback = "") {
    try {
      const v = localStorage.getItem(key);
      return v === null ? fallback : v;
    } catch {
      return fallback;
    }
  }

  function lsGetBool(key, fallback = false) {
    const v = lsGet(key, "");
    if (v === "") return fallback;
    return v === "true";
  }

  function lsSet(key, value) {
    try {
      localStorage.setItem(key, String(value));
    } catch {}
  }

  function loadPersisted() {
    selectedIndustry = lsGet(LS.industry, selectedIndustry);
    city = lsGet(LS.city, city);
    headless = lsGetBool(LS.headless, headless);

    lastSearchedIndustry = lsGet(LS.lastIndustry, "");
    lastSearchedCity = lsGet(LS.lastCity, "");
    setSearchDone(lsGetBool(LS.searchDone, false));
    setSearchAborted(lsGetBool(LS.searchAborted, false));
  }

  function persistInputs() {
    lsSet(LS.industry, selectedIndustry);
    lsSet(LS.city, city);
    lsSet(LS.headless, headless);
  }

  function setLastSearchedIndustry(industry) {
    lastSearchedIndustry = industry;
    lsSet(LS.lastIndustry, industry);
  }

  function setLastSearchedCity(cityValue) {
    lastSearchedCity = cityValue;
    lsSet(LS.lastCity, cityValue);
  }

  function setSearchDone(v) {
    searchDone = v;
    lsSet(LS.searchDone, v);
  }

  function setSearchAborted(v) {
    searchAborted = v;
    lsSet(LS.searchAborted, v);
  }

  onMount(() => {
    loadPersisted();
    if (searchDone) {
      loadScreenshots();
    }
  });

  $: persistInputs();

  async function searchForWebsites() {
    if (city !== "") {
      try {
        websiteUrls = await SearchForWebsites(
          city,
          selectedIndustry,
          !headless,
        );
        amountOfWebsites = websiteUrls.length;
        statusTextUrls = `${amountOfWebsites} Webseiten gefunden.`;
      } catch (error) {
        statusTextUrls = "";
        statusText = "Fehler beim Suchen der Webseiten: \n" + error.toString();
      }
    } else {
      statusTextUrls = "";
      statusText = "Stadt und Industrie müssen ausgefüllt sein!";
    }
  }

  async function takeScreenShotsOfWebsite(url) {
    try {
      await TakeScreenshotOfWebsite(url, !headless);
      statusText = "Screenshots von Webseiten machen...";
      return true;
    } catch (error) {
      statusTextScreenshots = "";
      statusText =
        `Fehler beim Erstellen eines Screenshots bei dieser Webseite:  ${url} \n` +
        error.toString();
      return false;
    }
  }

  function abortSearch() {
    setSearchAborted(true);
    setSearchDone(false);

    statusText = "";
    statusTextUrls = "";
    statusTextScreenshots = "";
    websiteUrls = [];
    amountOfWebsites = 0;
    screenshotProgress = 0;

    CancelCurrentJob();
  }

  async function main() {
    // new run
    setSearchAborted(false);
    setSearchDone(false);

    screenshotProgress = 0;
    amountSavedToExcel = 0;
    currentScreenshotIndex = 0;
    await ResetScreenshotsDir();

    setLastSearchedIndustry(selectedIndustry);
    setLastSearchedCity(city);

    statusText = "";
    statusTextScreenshots = "";
    statusTextUrls = "";
    websiteUrls = [];
    amountOfWebsites = 0;

    statusText = "URLs sammeln...";
    searchButtonActive = false;
    await searchForWebsites();

    if (searchAborted) {
      statusText = "";
      searchButtonActive = true;
      return;
    }

    statusText = "Screenshots von Webseiten machen...";
    for (let index = 0; index < websiteUrls.length; index++) {
      if (searchAborted) break;

      const websiteUrl = websiteUrls[index];
      let currentFiles = await GetScreenshotFilenames();
      statusTextScreenshots = `Bei Website: ${websiteUrl} \n${currentFiles.length} gespeichert.`;

      let success = await takeScreenShotsOfWebsite(websiteUrl);
      if (success) {
        screenshotProgress++;
      }
    }

    if (!searchAborted) {
      let finalFiles = await GetScreenshotFilenames();
      statusTextScreenshots = `${finalFiles.length} Screenshots von ${websiteUrls.length} Webseiten gemacht!`;
      setSearchDone(true);
      loadScreenshots();
    } else {
      statusTextScreenshots = "Suche abgebrochen.";
      setSearchDone(false);
    }

    statusText = "";
    searchButtonActive = true;
  }

  async function loadScreenshots() {
    try {
      screenshotFiles = await GetScreenshotFilenames();
      if (screenshotFiles.length > 0) {
        let savedIndex = await LoadImageIndex();
        if (savedIndex < 0 || savedIndex >= screenshotFiles.length) {
          savedIndex = 0;
        }
        currentScreenshotIndex = savedIndex;
        await updateScreenshotDisplay();
      }
    } catch (error) {
      console.error("Error loading screenshots:", error);
    }
  }

  async function updateScreenshotDisplay() {
    if (
      screenshotFiles.length > 0 &&
      currentScreenshotIndex >= 0 &&
      currentScreenshotIndex < screenshotFiles.length
    ) {
      try {
        currentScreenshotBase64 = await GetScreenshotBase64(
          screenshotFiles[currentScreenshotIndex],
        );
        // Persist current index to backend file
        await SaveImageIndex(currentScreenshotIndex);
      } catch (error) {
        console.error("Error loading screenshot base64:", error);
      }
    }
  }

  function nextScreenshot() {
    if (currentScreenshotIndex < screenshotFiles.length - 1) {
      currentScreenshotIndex++;
      updateScreenshotDisplay();
    }
  }

  function prevScreenshot() {
    if (currentScreenshotIndex > 0) {
      currentScreenshotIndex--;
      updateScreenshotDisplay();
    }
  }

  async function saveToExcel() {
    if (screenshotFiles.length > 0) {
      const filename = screenshotFiles[currentScreenshotIndex];
      const url = filename.replace(".png", "");
      try {
        await WriteBusinessUrlToExcelDatabase(url);
        amountSavedToExcel++;
        // Optional: you can show a subtle alert or log, but it might get spammy if pressing enter fast
        console.log(`Gespeichert: ${url}`);
      } catch (error) {
        alert(`Fehler beim Speichern: ${error}`);
      }
    }
  }

  function handleKeydown(event) {
    if (event.repeat) return;

    if (currentTab === "analyse_screenshots" && screenshotFiles.length > 0) {
      if (event.key === "ArrowRight") {
        nextScreenshot();
      } else if (event.key === "ArrowLeft") {
        prevScreenshot();
      } else if (event.key === "Enter") {
        saveToExcel();
        nextScreenshot();
      }
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<Tabs
  tabs={[
    { id: "website_search", label: "Webseiten suchen" },
    { id: "analyse_screenshots", label: "Bilder analysieren" },
  ]}
  canOpenAnalyse={searchDone}
  bind:active={currentTab}
>
  <div class="page" slot="website_search">
    <section class="card">
      <div class="form">
        <div class="grid">
          <div class="field">
            <label class="field-label" for="industries-select"
              >🏭 Industrie:</label
            >
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
            <label for="industries-select"
              >⏳ Zuletzt gesucht: {lastSearchedIndustry}</label
            >
          </div>

          <div class="field">
            <label class="field-label" for="city-select">🏙️ Stadt:</label>
            <select
              bind:value={city}
              required
              disabled={!searchButtonActive}
              name="cities"
              id="city-select"
            >
              {#each availableCities as c}
                <option value={c}>{c}</option>
              {/each}
            </select>
            <label for="city-select"
              >⏳ Zuletzt gesucht: {lastSearchedCity}</label
            >
          </div>
        </div>

        <div class="controls">
          <button
            class="btn btn--primary"
            on:click={main}
            disabled={!searchButtonActive}
          >
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

  <div class="analyse-page" slot="analyse_screenshots">
    <div class="analyse-container">
      <section class="card image-card">
        <div class="image-viewer">
          {#if screenshotFiles.length > 0 && currentScreenshotBase64}
            <img
              src={`data:image/png;base64,${currentScreenshotBase64}`}
              alt="Screenshot"
            />
          {:else}
            <div class="placeholder-text">Keine Screenshots gefunden.</div>
          {/if}
        </div>
      </section>

      <section class="controls-card">
        <div class="viewer-controls">
          <span class="label left">
            {screenshotFiles.length > 0
              ? `${currentScreenshotIndex + 1}/${screenshotFiles.length}`
              : "0/0"}
          </span>

          <span class="label left">
            {screenshotFiles.length > 0
              ? `${Math.round(((currentScreenshotIndex + 1) / screenshotFiles.length) * 100)}%`
              : "0%"}
          </span>

          <div class="buttons">
            <button
              class="btn btn--primary btn--icon"
              on:click={prevScreenshot}
              disabled={currentScreenshotIndex === 0 ||
                screenshotFiles.length === 0}
            >
              ←
            </button>

            <button
              class="btn btn--primary btn--icon"
              on:click={nextScreenshot}
              disabled={currentScreenshotIndex === screenshotFiles.length - 1 ||
                screenshotFiles.length === 0}
            >
              →
            </button>

            <button
              class="btn btn--primary"
              on:click={saveToExcel}
              disabled={screenshotFiles.length === 0}
            >
              💾 In Excel speichern
            </button>
            <span
              class="label left"
              style="display:flex; align-items:center; margin-left: 8px;"
            >
              Gespeichert: {amountSavedToExcel}
            </span>
          </div>

          <span class="label right url-label">
            {screenshotFiles.length > 0
              ? screenshotFiles[currentScreenshotIndex].replace(".png", "")
              : "---"}
          </span>
        </div>
      </section>
    </div>
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
    transition: all 0.2s ease;
  }

  .btn:hover:not(:disabled) {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
    border-color: rgba(0, 0, 0, 0.25);
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
    .viewer-controls {
      flex-direction: column;
      gap: 12px;
    }
    .label.right {
      margin-left: 0;
      text-align: center;
    }
  }

  .analyse-page {
    width: 100%;
    display: flex;
    justify-content: center;
    padding-bottom: 2rem;
  }

  .analyse-container {
    width: min(900px, 95vw);
    height: 80vh;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .image-card {
    flex: 1;
    display: flex;
    margin: 0;
    padding: 16px;
    /* The global .card adds min(760px), but .analyse-container bounds it here */
    width: 100%;
  }

  .controls-card {
    width: 100%;
    margin: 0;
    padding: 0;
    display: flex;
    justify-content: center;
  }

  .image-viewer {
    flex: 1;
    border: 1px solid rgba(0, 0, 0, 0.12);
    border-radius: 8px;
    overflow: hidden;
    display: flex;
    justify-content: center;
    align-items: center;
    background: #f5f5f5;
  }

  .image-viewer img {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
  }

  .placeholder-text {
    color: #888;
    font-size: 1.2rem;
  }

  .viewer-controls {
    display: flex;
    width: 100%;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 0;
    flex-wrap: wrap;
    border-top: none;
  }

  .buttons {
    display: flex;
    gap: 8px;
    justify-content: center;
  }

  .btn--icon {
    width: 42px;
    padding: 0;
    font-size: 1.2rem;
  }

  .label {
    font-size: 14px;
    font-weight: 600;
    color: #444;
    white-space: nowrap;
  }

  .url-label {
    flex: 1;
    text-align: right;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>
