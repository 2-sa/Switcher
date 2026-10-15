<script lang="ts">
  import { onMount } from "svelte";
  import { route, previousPage, appBarTitle } from "../stores/nav";
  import { flip } from "svelte/animate";
  import { t } from "../stores/i18n";
  import { DUR, fadeMotion, fadeUp, flipMotion } from "../lib/animation";
  import * as PlatformService from "../../bindings/TcNo-Acc-Switcher/internal/platform/platformservice.js";
  import "../styles/Settings.scss";
  import { homeScreenData } from "../stores/homeScreenData";

  let allNames: string[] = [];
  let disabled: string[] = [];
  let loadError: string | null = null;
  let loading = true;
  let saving = false;
  let saveError: string | null = null;

  $: appBarTitle.set($t("Title_Platforms_Settings"));
  $: disabledSorted = [...disabled].sort((a, b) =>
    a.localeCompare(b, undefined, { sensitivity: "base" }),
  );
  $: enabledSorted = allNames
    .filter((n) => !disabled.includes(n))
    .sort((a, b) => a.localeCompare(b, undefined, { sensitivity: "base" }));
  $: anyEnabled = enabledSorted.length > 0;

  async function refresh(): Promise<void> {
    loadError = null;
    try {
      const s = await PlatformService.GetStartup();
      homeScreenData.set(s);
      if (s.platformsFileMissing) {
        allNames = [];
        disabled = [];
        return;
      }
      allNames = s.allPlatformNames ?? [];
      disabled = [...(s.disabledPlatformNames ?? [])];
    } catch (e) {
      loadError = e instanceof Error ? e.message : String(e);
    } finally {
      loading = false;
    }
  }

  async function onPlatformChange(e: Event, item: string): Promise<void> {
    const el = e.currentTarget as HTMLInputElement;
    const enabled = el.checked;
    // Keep the checkbox consistent with saved state until persistence succeeds.
    el.checked = !disabled.includes(item);
    if (saving) return;
    saving = true;
    saveError = null;
    const next = enabled ? disabled.filter((x) => x !== item) : [...new Set([...disabled, item])];
    try {
      await PlatformService.SetDisabledPlatforms(next);
      disabled = next;
      await refresh();
    } catch (e) {
      saveError = $t("Shell_PlatformSaveError");
    } finally {
      saving = false;
    }
  }

  function saveAndClose(): void {
    route.set({ page: "home" });
  }

  onMount(() => {
    previousPage.set({ page: "home" });
    route.set({ page: "manage-platforms" });
    void refresh();
  });
</script>

<div class="main-content main-spacing manage-platforms" aria-busy={loading || saving}>
  <h1 class="SettingsHeader">{$t("Settings_Header_ExtraPlatforms")}</h1>

  {#if saveError}
    <p class="settings-err" role="alert">{saveError}</p>
  {/if}
  {#if loading}
    <p role="status">{$t("Button_Loading")}</p>
  {:else if loadError}
    <p class="settings-err" role="alert">{loadError}</p>
  {:else if !allNames.length}
    <span>{$t("Settings_NoPlatforms")}</span>
  {:else}
    {#if disabledSorted.length}
      <h2 class="SettingsHeader">{$t("Settings_ExtraPlatformsDisabled")}</h2>
      {#if !anyEnabled}
        <span>{$t("Settings_NoPlatforms")}</span>
      {/if}
      <div class="rowSetting platformsCheckboxes">
        {#each disabledSorted as item (item)}
          <div
            class="form-check mb-2"
            animate:flip={flipMotion()}
            in:fadeUp={{ y: 4, duration: DUR.fast }}
            out:fadeMotion={{ duration: DUR.instant }}
          >
            <input
              class="form-check-input"
              type="checkbox"
              id={"dis-" + item}
              disabled={saving}
              on:change={(e) => { void onPlatformChange(e, item); }}
            />
            <label class="form-check-label" for={"dis-" + item}></label>
            <label for={"dis-" + item}><bdi>{item}</bdi></label>
          </div>
        {/each}
      </div>
    {/if}

    {#if enabledSorted.length}
      <h2 class="SettingsHeader">{$t("Settings_ExtraPlatformsEnabled")}</h2>
      <div class="rowSetting platformsCheckboxes">
        {#each enabledSorted as item (item)}
          <div
            class="form-check mb-2"
            animate:flip={flipMotion()}
            in:fadeUp={{ y: 4, duration: DUR.fast }}
            out:fadeMotion={{ duration: DUR.instant }}
          >
            <input
              class="form-check-input"
              type="checkbox"
              id={"en-" + item}
              checked
              disabled={saving}
              on:change={(e) => { void onPlatformChange(e, item); }}
            />
            <label class="form-check-label" for={"en-" + item}></label>
            <label for={"en-" + item}><bdi>{item}</bdi></label>
          </div>
        {/each}
      </div>
    {/if}

    <div class="buttoncol col_close">
      <button type="button" class="btn_close" disabled={saving} on:click={saveAndClose}>
        <span>{$t("Button_Close")}</span>
      </button>
    </div>
  {/if}
</div>

<style lang="scss">
  .manage-platforms .platformsCheckboxes {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 8px;
    margin-block-end: 24px;
  }
  .manage-platforms .form-check {
    display: flex;
    align-items: center;
    gap: 8px;
    min-inline-size: 0;
    padding: 10px;
    border: 1px solid var(--shell-line);
    border-radius: 8px;
    background: var(--shell-raised);
  }
  .manage-platforms .form-check label:last-child {
    flex: 1;
    min-inline-size: 0;
    overflow-wrap: anywhere;
    cursor: pointer;
  }
  .settings-err {
    color: var(--red);
    padding: 0.5rem;
  }
</style>
