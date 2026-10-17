# Switcher — Interface Redesign

Status: redesign, Windows CI, and fork release publication verified; native visual/interaction and installed updater acceptance are pending.
Reviewed: 2026-10-16.

## [DISPLAY NAME — 2026-10-17]

- Follow-up: removed the redundant caption wordmark; Home's caption shows its translated page name, and translated page titles omit the product prefix. BearSwitcher remains the sidebar brand. Window controls keep their physical right-hand placement.

- User approved BearSwitcher as the visible name only. Updated interface headings, caption branding, tray labels, updater window, notifications, presence text, CLI heading, translated display strings, and Windows ProductName/FileDescription. Adjusted the sidebar wordmark size to fit the existing rail.
- Retained executable/build names, module/import names, bindings, application identifiers, protocol/IPC names, update URLs/signing/assets, settings/data paths, shortcut filenames, legal attribution, and author links. Translation text that names the real legacy data folder remains accurate.
- These are display changes only; native rendering and CI build validation are pending.

## [INSTALLER — 2026-10-17]

- User requires an installation wizard as the main download. Added a fork-specific NSIS installer with Arabic/English, directory selection, Start menu shortcuts, an optional desktop shortcut, WebView2 prerequisite detection, and per-user Windows uninstall registration. Default: LocalAppData/Programs/BearSwitcher, writable by the current user for signed in-app updates.
- Reuses the existing executable name and AppData settings/accounts. The fork installer does not enable reporting or copy the upstream installer's recursive data deletion. Uninstall deletes only installed files and its own registration/shortcuts; user data and unknown files are preserved.
- Packaging tools and smoke installation run only on GitHub. CI must exercise silent install, reinstallation/upgrade, uninstallation, hashes, shortcuts/registry, and preservation markers before signing and publishing. Both setup and updater payload receive detached Ed25519 signatures and SHA-256 checksums; these are not Windows Authenticode certificates.
- Official NSIS changelog identifies stable 3.13 (2026-09-27). The existing Wails beta.16 tool downloads Microsoft's signed Evergreen bootstrapper; its Authenticode signature is verified before packaging. CI results and interactive Arabic wizard acceptance remain pending.

## [PRIVACY]

- Branch verification: [run 38843086002](https://github.com/2-sa/Switcher/actions/runs/38843086002) passed on `996c6cfde4fb02ca90b91b1dbc0d3a7235c00e71`: Svelte/TypeScript, frontend tests/build, Go tests (including privacy regressions), and Windows build. Publication was intentionally skipped on the test branch.

- Removed statistics, crash, stability-rating, and feedback uploads to the original developer. Native compatibility methods cannot send data; reporting settings are forced off on load/save/batch updates, including legacy settings migration. Crash dumps remain local and are no longer associated with a statistics UUID.
- Local usage statistics default to off. Existing explicit local collection preferences remain supported. Sharing/automatic crash upload controls and the crash submission prompt were removed. Feedback opens the fork's GitHub issues page without prefilled account data, identifiers, or logs.
- Shared HTTP transport rejects `tcno.co` and its subdomains, including redirected requests and old/custom GameStats catalogs. Upstream-only Apex/CS fallback proxies, Steam catalog mirrors, and dynamic translator lists are consequently unavailable; cached data and independent platform sources remain usable. Intentional author/credit links in the external browser are unchanged.
- `steamguard.tcno.co` remains the existing local WebAuthn relying-party identifier; it is not an HTTP data upload destination. Changing it would invalidate existing hardware-key credentials.
- Regression coverage checks blocked requests/redirects, unaffected GitHub/Steam requests, old true flags, native setters/batches, disabled feedback/statistics APIs, and retention of local crash dumps. Native runtime network capture remains unverified.

## [NAVIGATION FOLLOW-UP]

- Moved the application back action out of the LTR native caption into a labeled page control at inline start: right in Arabic and left in English. Home no longer displays a nonfunctional back arrow. Standalone Steam confirmations retain their dedicated layout. Modal blocking, history behavior, keyboard/controller handling, and the controller hint remain supported.
- The sidebar is the sole global settings/platform-management entry point. The account footer retains a visibly named platform-specific settings action, such as `Steam settings`, instead of another generic gear.
- A browser layout preview using the changed styles confirmed RTL back placement to the right of the title, LTR placement to its left, no horizontal overflow at the inspected desktop size, no home back action, and one global settings/manage entry each. This is a layout preview with sample markup, not a rendered Wails application; native visual acceptance remains pending.

## [FORK_RELEASES]

- Every successful push to `master` in `2-sa/Switcher` builds and publishes a Windows x64 release through `verify.yml`. Versions use `4.1.<workflow run number>`; re-running a workflow retains its version. Master commits have independent concurrency groups so a later push does not cancel an earlier release build.
- The version is embedded before building. Go/frontend tests, the Windows build, Ed25519 signature verification against `updater-key.pub`, and SHA-256 checksum verification must succeed before publication. Files are uploaded to a draft before it is published; reruns preserve an already-published release.
- `FORK_UPDATER_KEY` is stored only as a repository Actions secret. The repository contains its public key and a signing test fixture. The upstream SignPath workflow remains separate; fork executables have updater signatures but no Windows Authenticode certificate.
- Existing builds with the upstream public key need one manual download of the first fork-signed release; subsequent fork-signed builds can update through the application.
- Verified first automatic publication: [run 38840878349](https://github.com/2-sa/Switcher/actions/runs/38840878349) passed on `6da32dd758bde5143cb2e0dd194a953426675e67` and published non-draft [v4.1.8](https://github.com/2-sa/Switcher/releases/tag/v4.1.8). The executable, detached signature, and SHA256SUMS are present. Version preparation was also exercised in an isolated temporary directory (app `4.1.123`, Windows `4.1.123.0`), and the pinned fork signing fixture verified locally using Node crypto.

## [SCOPE & ASSUMPTIONS]

- Redesign the entire desktop interface: shell, platform selection, accounts, settings, Steam tools, dialogs, menus, notifications, loading, empty, and error states.
- Approved direction: quiet dark UI, restrained accents, Arabic by default, English supported.
- Approved language policy: Arabic for fresh installations and missing language settings; preserve explicitly saved language preferences. Keep existing locale support.
- Preserve switching, account ordering, search, shortcuts, gamepad and keyboard navigation, Steam Guard, privacy controls, themes, backgrounds, and native window behavior.
- Keep existing service contracts and route identities. No new account features, framework migration, telemetry, or dependency upgrades in this scope.
- Initial native acceptance target: Windows, matching this workspace. Other operating systems require their own runtime checks before claiming compatibility.

## [TECH_STACK]

Versions declared in this checkout, not verified installed toolchains:

| Technology | Declared version | Decision |
| --- | --- | --- |
| Go | 1.26.5 | Preserve backend |
| Wails / JS runtime | 3.0.0-beta.16 | Preserve matching versions; this is a prerelease |
| Svelte | 5.57.0 | Reuse components and stores |
| TypeScript | 5.9.3 | Preserve current compiler |
| Vite | 8.2.2 | Preserve build pipeline |
| Vitest | 5.0.0 | Reuse frontend tests |
| Sass | 1.103.1 | Reuse SCSS and theme variables |
| pnpm | 11.0.8 | Respect packageManager and lockfile |

Official version research: [npm Svelte](https://www.npmjs.com/package/svelte?activeTab=versions), [npm Vite](https://www.npmjs.com/package/vite?activeTab=versions), [Wails releases](https://github.com/wailsapp/wails/releases), [TypeScript registry](https://registry.npmjs.org/typescript/latest).

The TypeScript registry returned latest `7.0.2`; Wails releases exposed `v3.0.0-beta.20`, while the project repository identifies v3 as beta. Latest stable Svelte, Vite, Vitest, and Sass versions remain unverified: local npm registry requests failed DNS resolution, and browser registry access was incomplete. Do not infer latest versions or deprecation status from the checkout. No upgrades are required for this visual redesign; verify any new dependency and API before adopting it.

Local discovery found Node and pnpm on PATH; Go and wails3 were not found, and frontend dependencies/bindings are absent. No dependencies or toolchains have been installed. The user explicitly deferred dependency installation and builds to GitHub when commits are uploaded. Do not install locally. Commit and push to `2-sa/Switcher` are now explicitly authorized; CI results must be checked after delivery.

## [DESIGN SYSTEM]

- Palette: canvas `#15181E`, navigation `#191D25`, raised surface `#222833`, primary text `#EDF1F7`, secondary text `#AAB5C5`, action accent `#8BAEF2`. Define semantic error/success colors separately and accompany them with text/icons.
- Typography: use local Segoe UI for Arabic and Latin initially; 14px body, 12px metadata, 20–24px page headings. Do not download fonts. Allow Arabic line height and avoid forced letter spacing.
- Layout: compact navigation rail on the right in Arabic, left in English; adjacent content with a page heading, relevant actions, and search. Collapse navigation when space is constrained.
- Identity: recognizable platform artwork and account avatars carry the visual emphasis. Use clear active-account status, calm surfaces, restrained borders, and few shadows. Avoid decorative dashboards and invented statistics.
- Accounts: consistent avatar/name/status/action hierarchy; retain existing view and ordering options. Keep destructive actions in clearly labeled menus or confirmation dialogs.
- Shared controls: coherent buttons, inputs, toggles, tabs, menus, dialogs, toasts, and focus indicators. Use different corner radii by hierarchy, not identical cards everywhere.
- Respect saved themes/backgrounds through existing variables and theme persistence. Validate the new default and representative dark/light custom themes.
- Use logical CSS properties and bidi isolation for IDs, URLs, paths, codes, and technical values. Mirror navigation by locale; do not blindly mirror logos or native window controls.
- Support keyboard access, accessible names, text contrast of at least 4.5:1, reduced motion, and compact windows without clipped actions.

Arabic layout concept; English reverses navigation placement:

```text
|                 Native title bar                    |
|-----------------------------------------------------|
| Page title + actions             | Platforms        |
| Search / relevant filters        | Steam            |
|                                 | Other platforms  |
| Platform tiles or account list   |                  |
|                                 | Manage platforms |
| Contextual feedback             | Settings         |
```

## [SYSTEM_FLOW]

1. Launch → load saved settings → choose saved language or Arabic default → display platforms, loading state, or actionable error.
2. Select platform → show actual account data → identify selected/current account → invoke existing switch action → show busy state → report the real result.
3. Open account actions → retain existing operations and confirmation rules → reflect the service response; never fabricate success.
4. Open settings → search and edit existing options → persist through existing services → show saved state; language switching updates direction immediately.
5. Open Steam tools → retain games, confirmations, server picker, browser, and clearing routes → apply consistent chrome and controls.
6. Open Steam Guard → retain lock, authentication, enrollment, restore, and recovery states → display sensitive values only under existing security rules.

Each flow must pass with Arabic and English, keyboard navigation, and loading/error/empty states where applicable.

## [ARCHITECTURE]

| Area | Existing files / boundaries | Planned work |
| --- | --- | --- |
| Shell | `frontend/src/App.svelte`, `components/TitleBar.svelte`, `ActionBar.svelte`, `stores/nav.ts` | Restructure navigation and visual hierarchy using existing routes |
| Shared platform navigation | `components/Sidebar.svelte`, `lib/platformLaunch.ts`, `pages/Home.svelte` | Share the existing executable-resolution/manual-location flow; serialize requests and keep sidebar ordering synchronized |
| Platforms | `pages/Home.svelte`, `ManagePlatforms.svelte`, `styles/HomePlatforms.scss` | Redesign platform discovery and management |
| Accounts | `pages/Platform.svelte`, `PlatformSteam.svelte`, `components/PlatformAccountsBase.svelte`, account SCSS | Unify account presentation without changing switching logic |
| Settings | `pages/Settings.svelte`, `PlatformSettings.svelte`, `components/settings/`, `styles/Settings.scss` | Reorganize presentation; preserve fields and persistence |
| Steam tools | `pages/Steam*.svelte`, `components/SteamGames*.svelte` | Apply matching layouts and states |
| Overlays | `components/AppModal.svelte`, `modals/`, `ContextMenu*.svelte`, `Toast*.svelte`, lock/progress overlays | Unify surfaces, focus, and feedback |
| Styling | `styles/theme.scss`, `style.scss`, `UI.scss`, `rtl.scss`, `lib/theme/` | Semantic design tokens and targeted component styles |
| Language | `stores/i18n.ts`, `Resources/ar-SA.json`, `Resources/en-US.json` | Arabic copy, English parity, direction and fallback |
| Default language | `internal/platform/settings.go`, `service_settings_getters.go`, `service_settings_setters.go`, `settings_batch.go` | Align default/fallback policy; update affected tests |

This is an impact map, not authorization to rewrite every listed file. Reuse existing shared settings and modal components; extract a component only when repeated UI actually needs it. No new Shared/Core layer or parallel navigation state.

## [SAFE LOGGING]

Reuse existing logging, redaction, and error-to-toast paths. Record only actionable failures and essential state transitions with debug/info/warn/error levels; exclude account names, IDs, tokens, codes, and secrets. Do not log hover, render, or routine navigation events.

The inspected `internal/filelog` sink and `internal/app/logsink.go` use synchronous file writes behind a mutex. This redesign adds no new logging transport or diagnostic I/O: native failures use existing service paths and UI error feedback. Do not claim that the existing logger is asynchronous. Any future diagnostic transport must enqueue without awaiting disk/network I/O, use a bounded queue with an explicit overflow policy, and flush on orderly shutdown. A global logging rewrite is outside this visual change.

## [MILESTONES & VERIFIABLE GOALS]

| Milestone | Work | Acceptance |
| --- | --- | --- |
| 0 — Baseline | Inventory visible actions, screenshot current UI, establish checks, inspect logger and language migration | Baseline results recorded; all existing actions mapped; logging/default policy resolved |
| 1 — Visual slice | Implement shell, design tokens, platform and account layouts | Review real Arabic/English screenshots; navigation and a real account-switch flow work |
| 2 — Complete interface | Apply system to settings, tools, dialogs, menus, and all overlays | Every existing route and dialog category reviewed; no old layout left untracked |
| 3 — Language and regression | Arabic defaults, translation parity, RTL/LTR, theme and accessibility checks | Saved preference preserved; fresh profile Arabic; English switch persists; focus and technical strings correct |
| 4 — Acceptance | Run relevant tests/builds and inspect native window flows | Tests/build results recorded; no new regression; native visual review complete; pending section honestly resolved |

Validation commands: `pnpm --dir frontend run check`, `pnpm --dir frontend run test`, `pnpm --dir frontend run build`, `go test ./...`, and `wails3 task build` when the toolchain is available. Record baseline failures separately. Use red/green regression tests for changed behavior; use visual and interaction checks for presentation rather than tests that duplicate CSS.

## [IMPLEMENTED & VERIFIED]

- Added a responsive, locale-aware navigation sidebar using existing routes and platform artwork. Standalone confirmations/browser windows retain their dedicated layout.
- Added page headings/search and moved Steam accounts/games controls from the native title bar into a visible page toolbar. Removed the superseded title-bar styles and toggle helper.
- Applied a slate/blue default palette, local Segoe UI typography, platform/account cards, settings groups, controls, action bars, search, context menus, dialog shells, and toasts through existing theme variables.
- Preserved explicit language settings; aligned missing-language defaults in the frontend and backend. Added regression cases for fresh settings and saved English/Arabic.
- Isolated account/platform names and technical input values; corrected logical toast/update/status positioning, modal margins, directional transitions, and the back arrow. Native window controls retain LTR physical placement.
- Shared native platform-resolution logic between Home and Sidebar, retaining cancellation/manual-path confirmation and shortcut feedback. Prevented overlapping navigation requests.
- Added `.github/workflows/verify.yml`: branch pushes, pull requests, and manual runs install pinned tools on Windows, generate bindings, check Svelte/TypeScript, run frontend/Go tests, build, and upload an unsigned test executable. Existing tag/signing release workflow is preserved.

Local checks performed without installing dependencies:

- Locale selection: seven Node-executed cases passed; Arabic default, preserved saved languages, and fallback behavior. A failing check was observed before the helper was added, followed by passing checks.
- Platform navigation: seven scenarios passed against the actual TypeScript helper using Node's module runner with native I/O stubs: resolved, cancelled, manually located, confirmation error, shortcut feedback, overlapping request, and unavailable platform.
- New Arabic/English shell keys and placeholder parity passed. Default text contrast checks passed: primary 15.68:1, secondary 7.13:1, accent 7.59:1 on their intended surfaces.
- Twenty-four changed TypeScript scripts (including Svelte script blocks) parsed with Node's TypeScript parser. JSON dictionaries parsed. This does not compile Svelte templates or perform semantic TypeScript checking.
- `git diff --check` and whitespace checks across changed/untracked source files passed.
- GitHub Verify run [38722942797](https://github.com/2-sa/Switcher/actions/runs/38722942797) passed on commit `69877c1ed73054c7ff301b943b6530434b1e3fde`: generated bindings, Svelte/TypeScript checking, Vitest, frontend production build, full Go suite, and Windows application build. The run took 4m41s. Artifact `switcher-windows-unsigned` is available as [artifact 11966333517](https://github.com/2-sa/Switcher/actions/runs/38722942797/artifacts/11966333517); it is a test executable, not a signed release. The user later downloaded and ran this test executable; baseline native UI inspection is recorded below.
- Reviewed a separate browser design preview with sample data: Arabic/English direction, platform/account/settings layouts, full and compact sidebar, and no document horizontal overflow at 600px width. Captured Arabic and English images. This preview reuses the palette/layout CSS but is not the running Wails application and does not prove native integrations.

Preview artifacts are outside the repository under `C:/Users/PC/.codex/visualizations/2026/10/15/01a13efc-f407-7740-b711-978244d25d3e/`: `redesign-preview.html`, `redesign-ar.jpg`, and `redesign-en.jpg`. A temporary loopback preview runs at `http://127.0.0.1:8347/`.

## [NATIVE UI REPAIR — 2026-10-15]

Baseline defects reproduced in the running downloaded executable: clipped footer actions at 902px, an empty pinned drop zone shifting controls when shortcuts open, the shortcut-folder button outside its popup, left-aligned Arabic settings headings, and hard-to-read account metadata. Source inspection additionally found unsynchronized sidebar data and unchecked/racing platform-setting writes.

- Footer actions now shrink to available width and wrap; status occupies its own row only when present. Login padding is compact and the footer keeps its required height.
- Shortcut menus and the Steam games account picker are constrained to content-pane bounds in both directions. The folder action stays inside the shortcut menu. Empty pinned targets appear only during dragging.
- Arabic headings, search results and context menus use logical alignment. Account lists have translated accessible names, visible current-account status, labeled last-use dates, and bidirectional date isolation. Search is visible on account pages; platform settings have a distinct accessible label.
- Platform management uses readable grid rows, a real loading state, serialized writes, saved-state preservation on failure, an inline error, and refreshed shared navigation data after loading.

Local validation without downloads/installations: nine dropdown regression cases passed against the actual TypeScript test source using a Node assertion shim (the new export was first verified missing); simulated the actual platform-change handler for failed writes, retry, concurrent writes, and successful commit; parsed ten changed TypeScript scripts; checked Arabic/English keys and placeholders; git whitespace check passed. Full Svelte compilation, Vitest and native build are pending GitHub verification. These local checks do not prove the repaired executable's rendered appearance.

## [RTL FOLLOW-UP — LOCAL REVIEW]

User feedback after commit 944bf850: the interface still contains widespread LTR text. Further commits/pushes are withheld pending review; the running executable does not automatically incorporate source edits.

- Audited every explicit text alignment declaration across Svelte components and SCSS. Removed remaining physical left/right text alignment in Steam Guard fields, unlock forms, account lists, restore screens, vault factors, confirmation rows/profile links, server table headers, file picker, feedback counter, action status, preview toast, captions, and the Steam2003 tooltip. Centered controls retain their intended alignment; technical input values remain isolated LTR.
- Converted search-icon gaps, factor-action spacing, file-tree indentation/badges, and context-menu gutters to logical properties. File labels use bdi; the path input explicitly preserves LTR.
- Arabic submenus now open toward the left and mirror their arrow. Keyboard expansion/back navigation follows locale direction in context menus and the file tree. Menu bounds handle overflow on either viewport edge.
- Local validation: the two new Arabic context-menu regression cases failed before the keyboard fix; all 16 actual test cases subsequently passed using a Node assertion/mock shim, without installing Vitest locally. Four file-tree arrow scenarios and three submenu-boundary simulations passed. Twelve changed TypeScript scripts parsed; whitespace checks passed. No physical left/right text-align declarations remain in application Svelte/SCSS sources.
- These are source and simulated-logic checks. Full Svelte compilation and the repaired native application's rendered RTL remain unverified. No additional commit, push, installation, or download was performed for this follow-up.

## [ORPHANS & PENDING]

### Updater restart investigation — 2026-10-16

- Native failure confirmed by the local Wails helper log: the update was staged and the helper started, but the parent process remained alive for 30 seconds; the helper aborted before replacing the executable. The user confirms both windows remain visible after Restart & Apply.
- The affected client reports 4.1.9, updating to 4.1.11. Its process and executable at the recorded download path were no longer available during investigation. The blocked shutdown hook has not been identified; no restart fix is claimed.
- Inspected the pinned Wails beta.16 restart/helper/host implementations and application shutdown hooks. An isolated copy of the already downloaded 4.1.11 executable was used for diagnostics; the failed button flow was not reproduced. Temporary test processes were stopped. No toolchain was installed and no release was published for this investigation.
- Next acceptance: reproduce using the affected client with debug logging, identify the blocked shutdown stage, add a regression test, and verify signed download → process exit → replacement → relaunch on Windows. Do not treat CI alone as confirmation.

### Bounded update shutdown — implementation

- Added an application-owned Wails updater host. Once Wails successfully starts its signed-update helper and calls Host.Quit, normal shutdown has a 20-second process-exit deadline, below the helper's 30-second wait. A blocked shutdown exits the process so the helper can replace and relaunch it. Ordinary exits and failed helper spawns do not arm this deadline. This is a bounded fallback; the original blocked shutdown hook remains unidentified.
- Ready notifications/replays and Restart & Apply use one guarded restart path. Repeated requests cannot launch competing helpers; helper-start errors are logged and allow retry.
- Added Go regression cases for blocked shutdown, asynchronous exit, duplicate restart requests, spawn-failure retry, and a real blocked child process. Local Go tools remain unavailable; CI validation and native signed-update acceptance are pending. Download the repaired release manually once because code in a new executable cannot repair the already-running old updater.

No known newly orphaned source remains. GitHub checks, build, signing verification, and release publication passed in run `38840878349`. A visual retest of the repaired executable and an installed update/download/restart flow remain pending; CI and publication do not establish those runtime results.

- Resolved CI failures: first run `38721926769` used `C:/Users/RUNNER~1/...` temporary paths, causing existing shortcut/backup path comparisons to fail. Initial correction `44adb716` used a runner context at unsupported job-env scope; `69877c1e` instead exports TEMP/TMP from RUNNER_TEMP through GITHUB_ENV in a startup step. The complete Go suite and executable build then passed in run `38722942797`. No protection checks or failing tests were disabled.
- Visually inspect the real Windows build: all routes, Steam games/tools, Guard enrollment/unlock/restore, error/loading/empty states, search/context menus, dialogs, drag ordering, keyboard/gamepad controls, and actual account switching.
- Validate saved custom backgrounds and representative dark/light themes in the real application; the sample preview covers the new default palette only.
- Baseline native Windows inspection now covers Home, Riot accounts, Steam accounts, shortcut expansion/Escape, settings search/language dropdown, and platform management at 902 x 598. Other OS runtime behavior is unverified.
- Latest-stable/deprecation verification remains incomplete for several upstream packages. No new package or upgrade was introduced; do not treat this redesign as a dependency modernization.
