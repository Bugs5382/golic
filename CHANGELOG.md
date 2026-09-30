# Golic Change Log

## v1.2.0 - 2026-09-30

### What Changed 👀

#### 🚀 Features

- feat(license-file): optionally write and check the LICENSE file @Bugs5382 (#53)

#### 🐛 Bug Fixes

- fix(replace): swap only golic's own header and keep the file's other comments and spacing @Bugs5382 (#50)

### Extra

**Full Changelog**: https://github.com/Bugs5382/golic/compare/v1.1.0...v1.2.0

## v1.1.0 - 2026-09-30

### What Changed 👀

#### 🚀 Features

- feat(config): ignore generated and vendored files by default @Bugs5382 (#47)

#### 🐛 Bug Fixes

- fix(config): correct XML and shebang placement, match nested Dockerfiles, add missing file types @Bugs5382 (#45)

### Extra

**Full Changelog**: https://github.com/Bugs5382/golic/compare/v1.0.0...v1.1.0

## v1.0.1 - 2026-09-29

### What Changed 👀

#### 🐛 Bug Fixes

- fix(config): correct XML and shebang placement, match nested Dockerfiles, add missing file types @Bugs5382 (#45)

### Extra

**Full Changelog**: https://github.com/Bugs5382/golic/compare/v1.0.0...v1.0.1

## v1.0.0 - 2026-09-27

### What Changed 👀

#### 💥 Breaking Changes

- chore!: mark the public interface stable for v1.0.0 @Bugs5382 (#41)
- refactor!: internalize the library packages and fix version reporting @Bugs5382 (#37)

#### 🐛 Bug Fixes

- fix(cli): print help when golic runs without a command @Bugs5382 (#39)
- fix(cli): exit 1 when files need changes and 2 on errors @Bugs5382 (#38)

#### 📄 Documentation

- docs(readme): use the Taskfile commands instead of make @Bugs5382 (#33)
- docs(readme): apply the lite emoji treatment @Bugs5382 (#30)

#### 🧩 Dependency Updates

- chore(deps): bump golang.org/x/sys to v0.44.0 @Bugs5382 (#40)

### Extra

**Full Changelog**: https://github.com/Bugs5382/golic/compare/v0.3.0...v1.0.0

## v0.3.0 - 2026-06-13

### What Changed 👀

#### 🚀 Features

- feat: add a replace command to update license headers @Bugs5382 (#25)
- feat: build releases with GoReleaser @Bugs5382 (#24)
- feat: add built-in license-header rules for TypeScript sources @Bugs5382 (#23)

#### 📄 Documentation

- docs: document the replace command and built-in TypeScript rules @Bugs5382 (#27)

### Extra

**Full Changelog**: https://github.com/Bugs5382/golic/compare/v0.2.0...v0.3.0

## v0.2.0 - 2026-04-20

### What Changed 👀

#### 🚀 Features

- fix: invalidate v0.1.1 @Bugs5382 (#17)

#### 🐛 Bug Fixes

- fix: invalidate v0.1.1 @Bugs5382 (#17)

### Extra

**Full Changelog**: https://github.com/Bugs5382/golic/compare/v0.1.1...v0.2.0

## v0.1.1 - 2026-04-15

### What Changed 👀

#### 🐛 Bug Fixes

- fix: double golic.yaml @Bugs5382 (#14)

### Extra

**Full Changelog**: https://github.com/Bugs5382/golic/compare/v0.1.0...v0.1.1

## v0.1.0 - 2026-04-12

### What Changed 👀

#### 🚀 Features

- feat: first release @Bugs5382 (#2)

#### 🐛 Bug Fixes

- fix: job release and bump @Bugs5382 (#8, #9, #10)
- fix: commit to git @Bugs5382 (#6)
- feat: added release-drafter template @Bugs5382 (#4)
- fix: disable cobra completion @Bugs5382 (#12)

### Extra

**Full Changelog**: https://github.com/Bugs5382/golic/compare/...v0.1.0
