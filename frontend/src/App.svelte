<script>
  import { SearchForWebsites } from "../wailsjs/go/main/App"
  import { TakeScreenshotOfWebsite } from "../wailsjs/go/main/App"
  import { ResetScreenshotsDir } from "../wailsjs/go/main/App"

  import data from "./data.json"
  let availableIndustries = data[1].industries
  let selectedIndustry = "Autoreparaturen"

  let statusText = ""
  let statusTextScreenshots = ""
  let statusTextUrls = ""
  let websiteUrls = []

  let city = "Merseburg"
  let headless = false

  let amountOfWebsites = 0

  let searchButtonActive = true

  async function searchForWebsites() {
    if (city != "") {
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
      statusText = `Fehler beim Erstellen eines Screenshots bei dieser Webseite:  ${url} \n` + error.toString()
    }
  }

  let screenshotProgress = 0
  async function main() {
    await ResetScreenshotsDir()
    statusText = ""
    statusTextScreenshots = ""
    statusTextUrls = ""

    statusText = "URLs sammeln..."
    searchButtonActive = false
    await searchForWebsites()

    statusText = "Screenshots von Webseiten machen..."
    for (let index = 0; index < websiteUrls.length; index++) {
      const websiteUrl = websiteUrls[index];
      statusTextScreenshots = `Bei Website: ${websiteUrl} \n` + `${index}/${websiteUrls.length}`
      await takeScreenShotsOfWebsite(websiteUrl)
      screenshotProgress++
    }
    statusTextScreenshots = `${screenshotProgress} Screenshots von ${websiteUrls.length} Webseiten gemacht!`
    statusText = ""

    searchButtonActive = true
  }
</script>

<main>
  <div id="inputs-container" class="inputs-container">
    <div id="inputs">
      <!-- Row 1: labels -->
      <div class="top-labels">
        <label for="industries-select">Industrie:</label>
        <label for="city-input">Stadt:</label>
      </div>

      <!-- Row 2: industry select + city input -->
      <div class="top-fields">
        <div id="industry-select-container" class="field industry-field">
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
        </div>

        <div class="field city-field">
          <!-- remove placeholder if you want it visually empty like the mockup -->
          <input
            disabled={!searchButtonActive}
            id="city-input"
            bind:value={city}
          />
        </div>
      </div>

      <!-- Row 3: search button + checkbox on the right -->
      <div class="middle-row">
        <button class="search-button" on:click={main} disabled={!searchButtonActive}>
          <!-- optional icon -->
          <span class="btn-icon" aria-hidden="true">⌕</span>
          <span>Suchen</span>
        </button>

        <label id="show-browser-checkbox" class="checkbox-area">
          <span class="checkbox-text">Browser<br />anzeigen</span>
          <input
            type="checkbox"
            disabled={!searchButtonActive}
            bind:checked={headless}
          />
          <span class="checkbox-visual" aria-hidden="true"></span>
        </label>
      </div>

      <!-- Row 4: full-width cancel button -->
      <button class="cancel-button" disabled={searchButtonActive}>
        <!-- optional icon -->
        <span class="btn-icon" aria-hidden="true">✖</span>
        <span>Suche abbrechen</span>
      </button>
    </div>
  </div>

  <!-- keep these at the bottom -->
  <p style="white-space: pre-line;">{statusTextUrls}</p>
  <p style="white-space: pre-line;">{statusText}</p>
  <p style="white-space: pre-line;">{statusTextScreenshots}</p>
</main>

<style>
/* Layout */
main {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

/* Keeps the status paragraphs below the UI block */
#inputs-container {
  width: 100%;
}

#inputs {
  display: grid;
  grid-template-columns: minmax(240px, 1.1fr) minmax(180px, 0.9fr);
  grid-template-areas:
    "labels labels"
    "industry city"
    "search checkbox"
    "cancel cancel";
  gap: 1rem 1.25rem;
  align-items: end;
  max-width: 900px;
}

/* Row 1 labels */
.top-labels {
  grid-area: labels;
  display: grid;
  grid-template-columns: minmax(240px, 1.1fr) minmax(180px, 0.9fr);
  gap: 1.25rem;
}

.top-labels label {
  font-size: 1.1rem;
  font-weight: 600;
}

/* Row 2 fields */
.top-fields {
  grid-area: industry / industry / city / city; /* occupy row, but internal grid handles columns */
  display: grid;
  grid-template-columns: minmax(240px, 1.1fr) minmax(180px, 0.9fr);
  gap: 1.25rem;
}

.field {
  min-width: 0;
}

.industry-field {
  grid-column: 1;
}

.city-field {
  grid-column: 2;
}

#industry-select-container,
.city-field {
  width: 100%;
}

#industries-select,
#city-input {
  width: 100%;
  height: 3.3rem;
  border: 2px solid #2b2b2b;
  border-radius: 16px;
  padding: 0 0.9rem;
  font-size: 1rem;
  background: white;
  box-sizing: border-box;
}

/* Better-looking select arrow (optional) */
#industries-select {
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  padding-right: 2.8rem;
  background-image:
    linear-gradient(45deg, transparent 50%, #000 50%),
    linear-gradient(135deg, #000 50%, transparent 50%);
  background-position:
    calc(100% - 22px) calc(50% - 3px),
    calc(100% - 12px) calc(50% - 3px);
  background-size: 10px 10px, 10px 10px;
  background-repeat: no-repeat;
}

/* Row 3 */
.middle-row {
  grid-area: search / search / checkbox / checkbox;
  display: grid;
  grid-template-columns: minmax(240px, 1.1fr) auto;
  gap: 1.25rem;
  align-items: center;
}

.search-button {
  grid-column: 1;
}

.checkbox-area {
  grid-column: 2;
  display: grid;
  grid-template-columns: auto 56px;
  grid-template-rows: auto;
  align-items: center;
  column-gap: 0.8rem;
  justify-self: end;
  cursor: pointer;
  user-select: none;
}

.checkbox-text {
  font-size: 1rem;
  line-height: 1.15;
}

/* Hide native checkbox but keep it accessible */
.checkbox-area input[type="checkbox"] {
  position: absolute;
  opacity: 0;
  pointer-events: none;
}

.checkbox-visual {
  width: 56px;
  height: 56px;
  border: 2px solid #2b2b2b;
  border-radius: 14px;
  background: white;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  position: relative;
}

/* Checkmark */
.checkbox-area input[type="checkbox"]:checked + .checkbox-visual::after {
  content: "✓";
  font-size: 2rem;
  line-height: 1;
  font-weight: 700;
  color: #111;
  transform: translateY(-1px);
}

/* Row 4 */
.cancel-button {
  grid-area: cancel;
  width: 100%;
}

/* Buttons (shared) */
.search-button,
.cancel-button {
  height: 3.5rem;
  border: 2px solid #2b2b2b;
  border-radius: 18px;
  background: white;
  font-size: 1.05rem;
  padding: 0 1rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.6rem;
  cursor: pointer;
  box-sizing: border-box;
}

.btn-icon {
  font-size: 1.2rem;
  line-height: 1;
}

.search-button:disabled,
.cancel-button:disabled,
#industries-select:disabled,
#city-input:disabled,
.checkbox-area input[type="checkbox"]:disabled + .checkbox-visual {
  opacity: 0.65;
  cursor: not-allowed;
}

/* Optional focus styles */
#industries-select:focus,
#city-input:focus,
.search-button:focus-visible,
.cancel-button:focus-visible,
.checkbox-area:focus-within .checkbox-visual {
  outline: 2px solid #666;
  outline-offset: 2px;
}

/* Responsive: stack on smaller screens */
@media (max-width: 700px) {
  #inputs,
  .top-labels,
  .top-fields,
  .middle-row {
    grid-template-columns: 1fr;
  }

  #inputs {
    grid-template-areas:
      "labels"
      "industry"
      "city"
      "search"
      "checkbox"
      "cancel";
  }

  .top-labels {
    display: grid;
    gap: 0.6rem;
  }

  .top-fields {
    grid-area: auto;
    display: grid;
    gap: 0.75rem;
  }

  .industry-field,
  .city-field {
    grid-column: auto;
  }

  .middle-row {
    grid-area: auto;
    display: grid;
    gap: 0.75rem;
  }

  .checkbox-area {
    justify-self: start;
  }

  .checkbox-text br {
    display: none;
  }
}
</style>
