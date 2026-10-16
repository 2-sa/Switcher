<script lang="ts">
  import { onMount } from "svelte";
  import { route, navigateBackLikeButton } from "../stores/nav";
  import { isRtl, t } from "../stores/i18n";
  import { activeModal } from "../stores/modal";
  import { contextMenu } from "../stores/contextMenu";
  import { searchOverlayCtrl } from "../stores/searchOverlay";

  let popupOpen = false;
  $: visible = $route.page !== "home" && $route.page !== "steam-confirmations";
  $: disabled = !!$activeModal;
  $: promptHidden = disabled || !!$contextMenu || $searchOverlayCtrl.open || popupOpen;

  onMount(() => {
    const refresh = () => {
      popupOpen = !!document.querySelector(
        ".shortcutDropdown.open, .dropdown-menu.show, .custom-dropdown-menu.show, [role='dialog']",
      );
    };
    const onUiStateChange = () => requestAnimationFrame(refresh);
    refresh();
    for (const event of ["focusin", "keydown", "pointerdown"]) {
      window.addEventListener(event, onUiStateChange, true);
    }
    return () => {
      for (const event of ["focusin", "keydown", "pointerdown"]) {
        window.removeEventListener(event, onUiStateChange, true);
      }
    };
  });
</script>

{#if visible}
  <button
    type="button"
    class="workspace-back"
    class:workspace-back--controller-prompt={!promptHidden}
    aria-label={$t("Aria_WindowBack")}
    {disabled}
    on:click={navigateBackLikeButton}
  >
    <svg viewBox="0 0 24 24" aria-hidden="true" style:transform={$isRtl ? "scaleX(-1)" : undefined}>
      <path d="m14 6-6 6 6 6" />
    </svg>
    <span>{$t("Aria_WindowBack")}</span>
  </button>
{/if}

<style>
  button.workspace-back {
    position: relative;
    display: inline-flex;
    flex: 0 0 auto;
    align-items: center;
    justify-content: center;
    gap: 6px;
    min-block-size: 38px;
    margin: 0;
    padding: 6px 10px;
    border: 1px solid var(--shell-line);
    border-radius: 8px;
    background: var(--shell-raised);
    color: var(--whiteSecondary);
    font: inherit;
    cursor: pointer;
  }

  button.workspace-back:disabled {
    opacity: 0.45;
    cursor: not-allowed;
  }

  .workspace-back svg {
    inline-size: 18px;
    block-size: 18px;
    fill: none;
    stroke: currentColor;
    stroke-width: 2;
    stroke-linecap: round;
    stroke-linejoin: round;
  }

  :global(.input-modality-controller) .workspace-back--controller-prompt::after {
    content: var(--controller-back-prompt-glyph);
    position: absolute;
    inset-block-start: -5px;
    inset-inline-end: -5px;
    display: grid;
    place-items: center;
    inline-size: 1rem;
    block-size: 1rem;
    border: 2px solid var(--controller-back-prompt-border);
    border-radius: 50%;
    background: var(--controller-back-prompt-bg);
    color: var(--controller-back-prompt-fg);
    font-size: 0.66rem;
    line-height: 1;
    pointer-events: none;
  }

  @media (forced-colors: active) {
    :global(.input-modality-controller) .workspace-back--controller-prompt::after {
      background: ButtonFace;
      border-color: ButtonText;
      color: ButtonText;
    }
  }
</style>
