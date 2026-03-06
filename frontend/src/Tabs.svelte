<script>
  export let tabs = [
    { id: "website_search", label: "Webseiten suchen" },
    { id: "analyse_screenshots", label: "Bilder analysieren" },
  ];

  // allow switching to analyse tab only when true
  export let canOpenAnalyse = false;

  export let active = tabs[0].id;

  $: isDisabled = (tab) => {
    return tab.id === "analyse_screenshots" && !canOpenAnalyse;
  };

  function activate(tab) {
    if (!isDisabled(tab)) active = tab.id;
  }

  // If analyse becomes disabled while active, jump back
  $: if (active === "analyse_screenshots" && !canOpenAnalyse) {
    active = tabs[0].id;
  }
</script>

<nav class="tabbar" role="tablist" aria-label="Tabs">
  {#each tabs as tab}
    <button
      type="button"
      role="tab"
      aria-selected={active === tab.id}
      aria-disabled={isDisabled(tab)}
      disabled={isDisabled(tab)}
      class="tab"
      class:active={active === tab.id}
      on:click={() => activate(tab)}
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
    width: min(760px, 100%);
    margin: 18px auto 0 auto;
    display: flex;
    gap: 10px;
    padding: 10px;

    border-radius: 14px;
    border: 1px solid rgba(0, 0, 0, 0.12);
    background: rgba(255, 255, 255, 0.9);
    box-shadow: 0 12px 30px rgba(0, 0, 0, 0.08);
    box-sizing: border-box;
  }

  .tab {
    flex: 1 1 0;
    height: 42px;
    padding: 0 14px;

    border-radius: 10px;
    border: 1px solid rgba(0, 0, 0, 0.16);
    background: white;

    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;

    font-weight: 600;
    font-size: 13px;
    white-space: nowrap;
    user-select: none;

    transition:
      background 150ms ease,
      border-color 150ms ease,
      transform 120ms ease,
      box-shadow 150ms ease;
  }

  .tab:hover {
    box-shadow: 0 6px 14px rgba(0, 0, 0, 0.08);
    transform: translateY(-1px);
  }

  .tab:active {
    transform: translateY(0px);
    box-shadow: none;
  }

  .tab:focus-visible {
    outline: none;
    box-shadow: 0 0 0 3px rgba(46, 204, 113, 0.25);
  }

  .tab.active,
  .tab[aria-selected="true"] {
    border-color: black;
    background-color: lightgrey;
  }

  .tabcontent {
    padding: 24px 0 0 0;
    display: flex;
    justify-content: center;
    box-sizing: border-box;
  }

  @media (max-width: 640px) {
    .tabbar {
      width: min(900px, calc(100% - 32px));
      margin-top: 12px;
      padding: 8px;
      gap: 8px;
    }

    .tab {
      height: 44px;
      padding: 0 10px;
      font-size: 13px;
    }
  }

  .tab:disabled,
  .tab[aria-disabled="true"] {
    opacity: 0.55;
    cursor: not-allowed;
    box-shadow: none;
    transform: none;
  }
</style>
