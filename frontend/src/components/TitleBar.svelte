<script lang="ts">
    import { route } from '../stores/nav'
    import { t } from "../stores/i18n";
    import { onMount } from 'svelte'
    import { appBarTitle } from '../stores/nav'
    import { Events, Window } from "@wailsio/runtime";
    import { gameMode } from '../stores/osCapabilities'

    let minimised = false
    let maximised = false

    const common = Events.Types.Common

    async function refreshWindowState() {
        const [min, max] = await Promise.all([
            Window.IsMinimised(),
            Window.IsMaximised(),
        ])
        minimised = min
        maximised = max
    }

    onMount(() => {
        void refreshWindowState()

        const tracked = [
            common.WindowMinimise,
            common.WindowUnMinimise,
            common.WindowMaximise,
            common.WindowUnMaximise,
            common.WindowRestore,
            common.WindowFullscreen,
            common.WindowUnFullscreen,
        ] as const

        const unsubs = tracked.map((subscribedAs) =>
            Events.On(subscribedAs, (ev) => {
                void refreshWindowState()
            })
        )

        return () => {
            for (const off of unsubs) off()
        }
    })

    $: titleLabel = (() => {
        if ($route.page === "manage-platforms") {
            return $t("Title_Platforms_Settings")
        }
        if ($route.page === "preview-css") {
            return $t("Title_Settings_TestCss")
        }
        if ($route.page === "settings") {
            return $t("Title_Settings")
        }
        if ($route.page === "platform-settings") {
            return $t("Title_Template_Settings", { platformName: $route.platformName })
        }
        if ($route.page === "platform") {
            return $t("Title_AccountsList", { platform: $route.platformName })
        }
        return $appBarTitle
    })()

</script>

<header class="headerbar" dir="ltr">
    <span class="title-left">
        <span class="header_icon">BearSwitcher</span>
    </span>
    <span id="title-label" class="title-drag" dir="auto">
        {titleLabel}
    </span>
    <span class="window-controls" role="toolbar">
        <!-- Game Mode fills the screen and the only way out is to close it, so
             there is nothing to minimise into and nothing to restore to. -->
        {#if !$gameMode}
        <button type="button" class="win-btn win-btn-min" aria-label={$t("Aria_WindowMinimize")} on:click={() => Window.Minimise()}>
            <svg class="win-btn__glyph win-btn__glyph--min" viewBox="0 0 10 10" aria-hidden="true">
                <path d="M1 5h8" />
            </svg>
        </button>
        {#if maximised}
            <button type="button" class="win-btn win-btn-max" aria-label={$t("Aria_WindowRestore")} on:click={() => Window.Restore()}>
                <svg class="win-btn__glyph win-btn__glyph--restore" viewBox="0 0 10 10" aria-hidden="true">
                    <path d="M2.5 1.5h5v5h-5z" />
                    <path d="M4.5 3.5h4v5h-5v-4" />
                </svg>
            </button>
        {:else}
            <button type="button" class="win-btn win-btn-max" aria-label={$t("Aria_WindowMaximize")} on:click={() => Window.Maximise()}>
                <svg class="win-btn__glyph win-btn__glyph--max" viewBox="0 0 10 10" aria-hidden="true">
                    <path d="M1.5 1.5h7v7h-7z" />
                </svg>
            </button>
        {/if}
        {/if}
        <button type="button" class="win-btn win-btn-close" aria-label={$t("Aria_WindowClose")} on:click={() => Window.Close()}>
            <svg class="win-btn__glyph win-btn__glyph--close" viewBox="0 0 10 10" aria-hidden="true">
                <path d="M2 2l6 6" />
                <path d="M8 2L2 8" />
            </svg>
        </button>
    </span>
</header>

<style lang="scss">
    .headerbar {
        --webkit-app-region: drag;
        --wails-draggable: drag;
        -moz-user-select: none;
        -ms-user-select: none;
        -webkit-user-select: none;
        user-select: none;
        z-index: 5;
        background: var(--border-bar-bg);
        position: relative;
        height: 32px;
        min-height: 32px;
        width: 100%;
        -webkit-app-region: drag;
        color: var(--whiteSecondary);
        grid-column: 1;
        display: flex;
        justify-content: space-between;
        align-items: center;
        overflow: hidden;
        font-family: "Segoe UI", sans-serif;
        font-size: 12px;
        font-weight: 500;
    }
    .title-left {
        z-index: 1;
        height: 100%;
        display: flex;
        flex-direction: row;
        align-items: center;
    }
    .title-drag {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        max-width: calc(100% - 300px);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        pointer-events: none;
    }
    .header_icon {
        margin: 12px;
        display: block;
        font-size: 10px;
        font-weight: 700;
        line-height: 1;
        color: var(--whiteSecondary);
    }
    .window-controls {
        --wails-draggable: no-drag;
        z-index: 1;
        display: grid;
        /* Sized to however many buttons are rendered: Game Mode shows only
           close, and a fixed three-slot grid would strand it away from the
           corner. */
        grid-auto-flow: column;
        grid-auto-columns: 46px;
        top: 0;
        right: 0;
        height: 100%;
    }
    .win-btn {
        border-radius: 0;
        background: none;
        border: 0;
        margin: 0;
        display: flex;
        justify-content: center;
        align-items: center;
        width: 46px;
        height: 100%;

        &:hover {
            background: var(--window-control-hover-bg);
        }
    }
    .win-btn__glyph {
        width: 10px;
        height: 10px;
        display: block;
        overflow: visible;
        color: currentColor;
        fill: none;
        stroke: currentColor;
        stroke-width: 1.2;
        stroke-linecap: square;
        stroke-linejoin: miter;
        vector-effect: non-scaling-stroke;
        forced-color-adjust: auto;
    }
    .win-btn__glyph--min {
        stroke-linecap: butt;
    }
    .win-btn-close:hover {background: var(--window-close-hover);}
</style>
