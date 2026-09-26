# AGENTS.md - golic

Guide for AI agents working in this repository. Pair with `CLAUDE.md` (the working agreement and
hook-enforced rules). Keep this file current when the build, layout, or public API changes.

## What this is

golic is a CLI that injects, replaces and removes license headers in source files. It ships as a
binary (`go install github.com/Bugs5382/golic/cmd/golic@<version>` or the GoReleaser archives).
There is no library API: everything except `cmd/golic` lives under `internal/`.

## Using golic

The contract is the CLI: the `inject`, `replace`, `remove` and `version` commands, their flags, the
`.golic.yaml` and `.licignore` formats, the `{{copyright}}` placeholder and the exit codes. The
hub's shared `job-golic` runs `golic inject --dry -x` and relies on a non-zero exit when headers are
missing, so treat any change there as breaking.

## Layout

- `cmd/golic/` - `main`, the process entry point.
- `internal/commands/` - the cobra commands and their flags.
- `internal/impl/` - the file walk, header rendering, config merge, and the embedded default
  ruleset (`default.golic.yaml`).
- `internal/build/` - version reporting: linker values first, then `runtime/debug.ReadBuildInfo`.
- `internal/` - options, glob matching, the service runner; `internal/logging/` sets up zerolog.

## Build, test, lint

- Build: `task build` (stamps `internal/build.Version` and `Gitsha` through `-ldflags`), or
  `CGO_ENABLED=0 go build ./...`.
- Test: `task test`, or `go test -race ./...`.
- Lint: `task lint` (goimports, golangci-lint, yamllint, gitleaks).
- License headers: `task license-dry` checks, `task license` stamps (both need `task build` first).

## Conventions and gotchas

- See `CLAUDE.md` for the branch/commit/PR rules; they are enforced by the git hooks in
  `.claude/hooks` (run `bash .claude/hooks/install.sh` once per clone).
- golic builds itself for its license check (`job-go-lic.yaml`); it does not use the hub's shared
  `job-golic`.
- The Release Manager rewrites the Taskfile `VERSION` and `CHANGELOG.md` on `main`. Leave both to it.
