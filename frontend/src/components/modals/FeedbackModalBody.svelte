<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { openExternalUrl } from "../../lib/openExternalUrl";
  import { t } from "../../stores/i18n";

  export let mode: "issue" | "suggestion" = "issue";
  export let platform = "";

  const dispatch = createEventDispatcher<{ resolve: string | null }>();

  async function openIssues(): Promise<void> {
    if (await openExternalUrl("https://github.com/2-sa/Switcher/issues/new")) {
      dispatch("resolve", null);
    }
  }
</script>

<div class="modal-block">
  <p>{$t("Feedback_ForkPrivacy")}</p>
  {#if mode === "issue" && platform}
    <p><bdi>{platform}</bdi></p>
  {/if}
  <div class="modal-buttons">
    <button type="button" class="btnicontext" on:click={() => dispatch("resolve", null)}>
      {$t("Button_Cancel")}
    </button>
    <button type="button" class="btnicontext modal-primary" on:click={() => void openIssues()}>
      {$t("Feedback_OpenFork")}
    </button>
  </div>
</div>
