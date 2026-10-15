# Repository Guidelines

## Project Structure & Module Organization

TcNo Account Switcher is a Go/Wails v3 desktop application with a Svelte 5 and TypeScript frontend.

- Root Go files provide the entry point, filesystem services, and embedded catalogs.
- `internal/` contains backend packages: `platform/` manages platform configuration, `basic/` handles generic switching, and `steam/` and `steamguard/` implement Steam workflows.
- `frontend/src/` contains `components/`, `pages/`, `stores/`, shared `lib/` helpers, `styles/`, and translation `Resources/`. Static assets live in `frontend/public/`.
- `Platforms*.json` and `GameStats.json` define platform and statistics catalogs.
- `build/` holds OS-specific tasks and packaging assets; `docs/` contains security guidance and architecture decisions.

## Build, Test, and Development Commands

Use the Go version in `go.mod`, pnpm version in `frontend/package.json`, and matching Wails v3 CLI. Run these from the repository root:

- `pnpm --dir frontend install --frozen-lockfile`: install locked frontend dependencies.
- `wails3 task dev`: launch desktop development mode with Vite.
- `wails3 task build`: build the application for the host OS into `bin/`.
- `pnpm --dir frontend run check`: run Svelte and TypeScript checks.
- `pnpm --dir frontend run build`: build production frontend assets.
- `go test ./...`: run Go tests, including root-level tests.
- `pnpm --dir frontend run test`: run Vitest once; use `test:watch` for iteration.

The current Taskfile `test` task references a missing `cmd/` directory; use the direct test commands above.

## Coding Style & Naming Conventions

Format Go with `gofmt`; use idiomatic exported names and OS suffixes such as `_windows.go`. Match existing frontend style: two-space indentation, double-quoted TypeScript strings, PascalCase Svelte components, and camelCase helpers. No dedicated ESLint or Prettier script is configured. Regenerate Wails bindings through build tasks when backend APIs change.

## Testing Guidelines

Place Go tests beside implementation as `*_test.go`, and frontend tests as `src/**/*.test.ts`. Use Go's testing package and Vitest. Add regression coverage for changed behavior; no numeric coverage threshold is configured. Validate native integrations on the affected OS and visually check UI changes.

## Commit & Pull Request Guidelines

History mixes imperative subjects with `fix:` and `feat:` prefixes; use concise, descriptive subjects. Follow `CONTRIBUTING.md`: open or reference an existing issue, announce intent there, and keep PRs focused. Include the linked issue, behavior changes, validation results, and screenshots for visual changes.

## Security & Configuration

Read `docs/steam-guard.md` before changing authenticator workflows. Keep credentials, authenticator secrets, account data, and unredacted diagnostic logs out of commits and PR attachments.
