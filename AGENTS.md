# AGENTS.md - golic

Guide for AI agents working in this repository. Pair with `CLAUDE.md` (the working agreement and
hook-enforced rules). Keep this file current when the build, layout, or public API changes.

## What this is

golic is a CLI that injects, replaces and removes license headers in source files. It ships as a
binary (`go install github.com/Bugs5382/golic/cmd/golic@<version>` or the GoReleaser archives).
There is no library API: everything except `cmd/golic` lives under `internal/`.

## Using golic

The contract is the CLI: the `inject`, `replace`, `remove` and `version` commands, their flags
(including `--license-file` on `inject` and `replace`), the `.golic.yaml` and `.licignore` formats
(including the optional `licenseFile` key), the `{{copyright}}` placeholder and the exit codes. The
hub's shared `job-golic` runs `golic inject --dry -x` and relies on a non-zero exit when headers are
missing, so treat any change there as breaking.

## Layout

- `cmd/golic/` - `main`, the process entry point.
- `internal/commands/` - the cobra commands and their flags.
- `internal/impl/` - the file walk, header rendering, config merge, and the embedded default
  ruleset (`default.golic.yaml`). `ignore.go` decides which files are in scope (`.licignore` over
  the config `ignore:` list); `generated.go` skips files with a generated-code marker.
  `licensefile.go` writes the LICENSE file from the official texts embedded under `licenses/`,
  each downloaded verbatim from the source it names; a test pins their checksums.
- `internal/build/` - version reporting: linker values first, then `runtime/debug.ReadBuildInfo`.
- `internal/` - options, glob matching, the service runner; `internal/logging/` builds the logger on
  `github.com/Bugs5382/go-log`'s neutral `TraceLogger`, and every package logs through
  `logging.L()` with `golog.F` fields. Do not import `github.com/rs/zerolog` anywhere; it is only
  an indirect dependency of go-log. Use `logging.DebugEnabled()` to guard costly debug output.

## Build, test, lint

- Build: `task build` (stamps `internal/build.Version` and `Gitsha` through `-ldflags`), or
  `CGO_ENABLED=0 go build ./...`.
- Test: `task test`, or `go test -race ./...`.
- Lint: `task lint` (goimports, golangci-lint, yamllint, gitleaks).
- License headers: `task license-dry` checks, `task license` stamps (both need `task build` first).
  Both pass `--license-file`, so golic's own `LICENSE` stays the verbatim apache.org text.

## Conventions and gotchas

- See `CLAUDE.md` for the branch/commit/PR rules; they are enforced by the git hooks in
  `.claude/hooks` (run `bash .claude/hooks/install.sh` once per clone).
- Open every PR as a draft. CI skips drafts, so run the full checks locally, push once they pass, and mark the PR ready when the work is finished; see CLAUDE.md "CI and Actions minutes".
- golic finds its own header by the exact rendered text. Change a built-in licence text only by
  adding the old text to `superseded.golic.yaml`, so files stamped by older versions keep working.
  Never edit an existing rule's `prefix`/`suffix` in `default.golic.yaml`: every stamped file in
  every repo would get a second header. Adding rules or changing `under` is safe.
  `TestNoChurnOnFilesStampedByV100` replays v1.0.0 output for each rule and fails on any such
  change.
- Built-in `ignore:` patterns and the generated-file skip only take files out of scope. Adding a
  pattern is safe; removing one puts files back in scope in every repo that relies on it.
- golic builds itself for its license check (`job-go-lic.yaml`); it does not use the hub's shared
  `job-golic`.
- The Release Manager rewrites the Taskfile `VERSION` and `CHANGELOG.md` on `main`. Leave both to it.
