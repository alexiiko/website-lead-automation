<script>
  export let tabs = [
    {id: "website_search", label: "Webseiten suchen"},
    {id: "analyse_screenshots", label: "Bilder analysieren"}
  ]

  let active = tabs[0].id
</script>

<nav class="tabbar" role="tablist" aria-label="Tabs">
    {#each tabs as tab}
      <button
        type="button"
        role="tab"
        aria-selected={active === tab.id}
        class="tab"
        class:active={active === tab.id}
        on:click={() => (active = tab.id)}
      >
        {tab.label}
      </button>
    {/each}
</nav>

<section class="tabcontent">
  {#if active === "website_search"}
    <slot name="website_search" />
  {:else if active === "analyse_screenshots"}
    <slot name="analyse_screenshots" />
  {/if}
</section>

<style>
  .tabbar {
    width: min(900px, calc(100% - 32px));
    margin: 18px auto 0 auto; /* centered, top */
    display: grid;
    grid-auto-flow: column;
    grid-auto-columns: 1fr;

    background: #fff;
    border: 1px solid #1f1f1f;        /* thick outline like the image */
    border-radius: 999px;               /* pill shape */
    overflow: hidden;                   /* keep buttons clipped to pill */
    box-shadow: 1px 1px lightgrey;
  }

  /* Tab buttons */
  .tab {
    height: 48;                       /* big like the screenshot */
    padding: 0 28px;
    font-size: 22;                    /* adjust to taste */
    font-weight: 500;
    letter-spacing: 0.2px;

    border: 0;
    background: transparent;
    cursor: pointer;

    display: grid;
    place-items: center;
    user-select: none;
  }

  /* Divider between tabs */
  .tab + .tab {
    border-left: 10px solid #1f1f1f;
  }

  /* Active tab styling */
  .tab.active {
    background-color: lightgrey;
  }

  /* Content area */
  .tabcontent {
    flex: 1;
    padding: 24px;
    display: flex;
    justify-content: center;
  }

  /* Make it behave on smaller screens */
  @media (max-width: 700px) {
    .tab {
      height: 70px;
      font-size: 22px;
      padding: 0 16px;
    }
    .tabbar {
      border-width: 6px;
      border-bottom-width: 8px;
    }
    .tab + .tab {
      border-left-width: 6px;
    }
  }
</style>
