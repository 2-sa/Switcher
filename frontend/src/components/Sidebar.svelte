<script lang="ts">
  import { route, type Route } from "../stores/nav";
  import { homeScreenData } from "../stores/homeScreenData";
  import { t } from "../stores/i18n";
  import { platformIconFgHref } from "../lib/platformIcon";
  import { platformDisplayLabels } from "../lib/platformName";
  import { platformActionBusy } from "../stores/platformPage";
  import { securityStatus } from "../stores/security";
  import { activeModal } from "../stores/modal";
  import { openPlatformPage, platformNavigationBusy } from "../lib/platformLaunch";
  import { pushToast } from "../stores/toast";

  let navigating = false;

  $: platforms = $homeScreenData?.homePlatformOrder ?? [];
  $: labels = platformDisplayLabels(platforms);
  $: disabled = navigating || $platformNavigationBusy || $platformActionBusy.busy || $securityStatus.appLocked || !!$activeModal;
  $: activePlatform = "platformName" in $route
    ? $route.platformName
    : $route.page.startsWith("steam-") ? "Steam" : "";

  async function navigate(next: Route): Promise<void> {
    if (disabled) return;
    navigating = true;
    try {
      if (next.page === "platform") await openPlatformPage(next.platformName);
      else route.set(next);
    } catch (error) {
      pushToast({ type: "error", message: error instanceof Error ? error.message : String(error), duration: 8000 });
    } finally {
      navigating = false;
    }
  }
</script>

<nav class="sidebar" aria-label={$t("Shell_Navigation")}>
  <div class="sidebar__brand" dir="ltr">
    <svg viewBox="0 0 32 32" aria-hidden="true"><path d="M5 9h19l-4-4m7 18H8l4 4M24 9l-4 4M8 23l4-4" /></svg>
    <span>BearSwitcher</span>
  </div>
  <button class="sidebar__link" class:sidebar__link--active={$route.page === "home"}
    aria-current={$route.page === "home" ? "page" : undefined} {disabled}
    on:click={() => void navigate({ page: "home" })} title={$t("Shell_Platforms")} aria-label={$t("Shell_Platforms")}>
    <svg viewBox="0 0 24 24" aria-hidden="true"><rect x="3" y="3" width="7" height="7" rx="2" /><rect x="14" y="3" width="7" height="7" rx="2" /><rect x="3" y="14" width="7" height="7" rx="2" /><rect x="14" y="14" width="7" height="7" rx="2" /></svg>
    <span>{$t("Shell_Platforms")}</span>
  </button>
  <div class="sidebar__platforms" role="group" aria-label={$t("Shell_Platforms")}>
    {#each platforms as name (name)}
      <button class="sidebar__link" class:sidebar__link--active={activePlatform === name}
        aria-current={activePlatform === name ? "page" : undefined} {disabled}
        on:click={() => void navigate({ page: "platform", platformName: name })} title={labels.get(name) ?? name} aria-label={labels.get(name) ?? name}>
        <svg class="sidebar__platform-icon" viewBox="0 0 500 500" aria-hidden="true"><use href={platformIconFgHref(name)} /></svg>
        <bdi>{labels.get(name) ?? name}</bdi>
      </button>
    {/each}
  </div>
  <div class="sidebar__utilities">
    <button class="sidebar__link" class:sidebar__link--active={$route.page === "manage-platforms"}
      aria-current={$route.page === "manage-platforms" ? "page" : undefined} {disabled}
      on:click={() => void navigate({ page: "manage-platforms" })} title={$t("Button_ManagePlatforms")} aria-label={$t("Button_ManagePlatforms")}>
      <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 4v16M4 12h16" /></svg>
      <span>{$t("Button_ManagePlatforms")}</span>
    </button>
    <button class="sidebar__link" class:sidebar__link--active={$route.page === "settings"}
      aria-current={$route.page === "settings" ? "page" : undefined} {disabled}
      on:click={() => void navigate({ page: "settings" })} title={$t("Shell_Settings")} aria-label={$t("Shell_Settings")}>
      <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M4 7h16M4 17h16" /><circle cx="9" cy="7" r="3" /><circle cx="15" cy="17" r="3" /></svg>
      <span>{$t("Shell_Settings")}</span>
    </button>
  </div>
</nav>
