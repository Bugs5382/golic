# Golic 🚀

> 📜 **A declarative tool for injecting and managing license headers in source code.**

Golic automates the tedious task of ensuring every source file in your project has the correct license header. It’s built for developers who want a "set it and forget it" solution for compliance.

```bash
# Preview changes before applying
golic inject -t apache2 -c "2026 MyCompany ltd." --dry
```

## 📥 Installation

Install the binary using Go 1.26+:

```bash
go install github.com/Bugs5382/golic/cmd/golic@latest
golic version
```

## ⚙️ Configuration

Golic relies on two small files in your project root to define **what** to license and **how** to format it.

### 1\. `.licignore`

Determines which files Golic should touch. It uses standard `.gitignore` syntax.

**Pro Tip:** Use "inverse rules" to be safe. Deny everything by default, then allow specific files:

```bash
# .licignore

# 1. Ignore everything
*

# 2. Allow specific files/patterns
!Dockerfile*
!Makefile
!*.go

# 3. Allow subdirectories
!*/
```

### 2\. `.golic.yaml`

Contains license text and formatting rules. Golic merges your local file with its [embedded master configuration](internal/impl/default.golic.yaml) by default.

Comment rules for common file types are **built in**. You only need a local `rules:` entry for unusual file types or to override a built-in. Which files are processed is controlled separately by `.licignore`, so to stamp TypeScript sources you just allow them there (e.g. `!*.ts`).

| Comment style | Built-in file types |
| :--- | :--- |
| `/* … */` | `.go` (below `package`), `.java`, `.scala`, `.cs`, `.kt`, `.swift`, `.rs`, `.c`, `.h`, `.cpp`, `.proto`, `.js`, `.jsx`, `.mjs`, `.cjs`, `.ts`, `.tsx`, `.mts`, `.cts`, `.css`, `.scss`, `.less` |
| `# ` | `.yaml`, `.yml`, `.toml`, `.hcl`, `.tf`, `.tfvars`, `.graphql`, `.graphqls`, `.py`, `.sh`, `.bash`, `.zsh`, `.gitignore`, `.dockerignore`, `.helmignore`, `.licignore`, `Dockerfile*`, `Containerfile*`, `Makefile`, `.mk` |
| `-- ` | `.sql` |
| `<!-- … -->` | `.html`, `.xml`, `.svg` |
| `{{/* … */}}` | `.tpl` |

`Dockerfile*`, `Containerfile*` and `Makefile` match at any depth, so `build/Dockerfile` and `sub/Makefile` are covered too.

The header always goes below the lines a file needs first:

* 🐚 **Shebangs:** a `#!` first line in shell, Python, JavaScript and TypeScript files stays on line 1, whatever the interpreter path (`#!/usr/bin/env node`, `#!/usr/local/bin/bash`, `#!/bin/sh`).
* 📄 **XML declarations:** `<?xml … ?>` stays first in `.xml` and `.svg` files.
* 🐳 **Dockerfile parser directives:** `# syntax=`, `# escape=` and `# check=` stay above the header, so BuildKit still reads them.

JSON, Markdown and MDX have no built-in rule on purpose: JSON has no comments, and a header in Markdown shows up in the rendered page.

> [!WARNING]
> A header in a `.sql` file changes its checksum. Migration tools that checksum applied files, such as Atlas (`atlas.sum`) and Flyway, will report already-applied migrations as modified. Keep migration folders out of `.licignore`, or stamp them before they are first applied.

```yaml
# .golic.yaml 
golic:
  licenses:
    apacheX: |
      Copyright MyCompany
      License details go here...
      
  rules:
    "*.go.txt":
      prefix: "/*"
      suffix: "*/"
    "technitium/**/*.yaml":
      prefix: "{{/*"
      suffix: "*/}}"
    .mzm:
      prefix: "" # No indent/prefix; place text directly at top
```

A rule key is a file extension (`.go`) or a glob matched against the file path relative to the directory golic runs in (`"**/Dockerfile*"`). When several keys match a file, the longest key wins. A rule has three fields:

* `prefix`: starts every header line, or opens the block when `suffix` is set.
* `suffix`: closes a block comment. Leave it out for line comments.
* `under`: line patterns the header goes below, for example `"package *"` or `"#!/**"`. The first pattern that matches wins, and the lines right after it that match any pattern stay above the header too. Patterns starting with `#!` only match the first line.

> [!IMPORTANT]
> Golic recognises a header it wrote by its exact text. Changing the `prefix` or `suffix` of a rule you have already stamped with makes every file look unstamped, and the next run adds a second header. Changing `under` is safe.

#### Merging with the built-in rules (`mergeRules`)

`mergeRules` decides how the `rules:` in your `.golic.yaml` combine with the built-in rules. Licenses always merge, whatever it is set to.

* **`true` (the default):** your rules are added to the built-in set. A rule with the same key as a built-in replaces that built-in rule as a whole; fields are not merged, so an override of `.sh` that leaves out `under` no longer keeps the shebang first.
* **`false`:** your rules **replace the built-in set entirely**. Only the keys you list exist. Every other file type has no rule and is **skipped without a warning**, even when `.licignore` allows it. If you set `mergeRules: false` and list no rules at all, the built-in set is kept.

With `mergeRules: false` you have to list every rule you need. This config stamps Go and YAML files and nothing else; shell scripts, Dockerfiles and TypeScript are silently skipped:

```yaml
# .golic.yaml
golic:
  mergeRules: false
  rules:
    .go:
      prefix: "\n/*"
      suffix: "*/"
      under:
        - "package *"
    .yaml:
      prefix: "# "
```

To keep a type covered, copy its rule from the [embedded master configuration](internal/impl/default.golic.yaml) byte for byte, so files you already stamped are still recognised.

## 🛠️ Usage

### Injecting Licenses

Once your config is ready, run:

```bash
golic inject -t apacheX
```

> [!TIP]
> Always use the `--dry` flag first to see a preview of which files will be modified without actually changing them.

### Updating or Removing

**Update in one pass** with `replace` — it strips the existing license header and injects the configured one, so you can change the copyright holder, year, or license type without a separate remove step:

```bash
# Preview first, then apply
golic replace -c "2026 MyCompany ltd." -t apacheX --dry
golic replace -c "2026 MyCompany ltd." -t apacheX
```

`replace` is a no-op when a file already carries the exact target header, and honors `--dry` and `-x` like `inject`.

**Remove** headers entirely with:

```bash
golic remove -t apacheX
```

### CI/CD Integration

To fail a build if licenses are missing (e.g., a developer forgot to run Golic), add `-x` (`--modified-exit`) to a dry run:

```bash
golic inject --dry -x -t apache2
```

With `-x`, golic exits 1 when any file would change and lists the count, for example `1 file(s) need a license header`. Without `--dry`, `-x` still exits 1 after it modifies files.

### Exit Codes

| Code | Meaning |
| :--- | :--- |
| `0` | Clean: nothing needed a change, or `-x` was not set. |
| `1` | `-x` was set and files were modified, or in a dry run need a header added, replaced or removed. |
| `2` | A real error: an unknown flag, a missing template, a license key or config that cannot be found or parsed, or a file that cannot be read or written. |

CI jobs should treat any non-zero status as a failure. Scripts that need to tell "headers missing" apart from "golic could not run" can check for `1` and `2`.

## 🏗 Development

The project uses [Task](https://taskfile.dev) (`Taskfile.yaml`). Run `task --list` to see every target.

### Build

To compile the project locally, execute:

```bash
task build
```

The binary lands in `bin/golic-<os>-<arch>`.

To remove build artifacts and clean your workspace:

```bash
task clean
```

If you are **contributing** to this project, you must first initialize the linting environment:

```bash
task lint-init
```

This command installs all necessary dependencies and tools for code analysis.

Once initialized, you can analyze the codebase by running:

```bash
task lint
```

`task lint` stamps missing license headers first, using the binary from `task build`.

To verify only the project licenses without changing any files (exit 1 when a header is missing), use:

```bash
task license-dry
```

To add any missing headers, run `task license`.

### Test

To execute the unit testing suite, run:

```bash
task test
```

## 📋 Command Reference

| Command | Description |
| :--- | :--- |
| `inject` | Adds the license header to files that are missing it. |
| `replace` | Replaces an existing header with the configured one in a single pass. |
| `remove` | Removes the license header matching the config. |
| `version` | Prints the version, and the commit on a second line when known. |

**Common Flags:**

* `-t, --template` : The license key to use from the config, for example `apache2` or `mit`. Required.
* `-c, --copyright` : Set the holder/year (Default: `YYYY [Insert Company]`). YYYY defaults to the current year.
* `-d, --dry` : Report what would change without writing any files.
* `-x, --modified-exit` : Exit with status 1 when any file is modified, or would be in a dry run. Errors exit with 2.
* `-l, --licignore` : Path to the `.licignore` file (Default: `.licignore`).
* `-p, --config-path` : Path to a local config merged over the built-in rules and licenses (Default: `.golic.yaml`).
* `-v, --verbose` : Enable detailed trace logging.

Running `golic` with no command prints the help.

## 🔒 Stability

From v1.0.0 the golic CLI is stable. These change only in a new major version:

* the `inject`, `replace`, `remove` and `version` commands and their flags;
* the `.golic.yaml` and `.licignore` file formats;
* the `{{copyright}}` placeholder in license templates;
* the exit codes (`0` clean, `1` files need changes with `-x`, `2` errors).

Minor releases may add commands, flags, built-in licenses and comment rules without breaking existing configs. golic has no Go library API: everything outside `cmd/golic` lives under `internal/`.

## 🤝 Contributing

We welcome Pull Requests! Please follow these steps:

* ✅ **Validation:** Run `task lint` to verify code quality.
* 🧪 **Testing:** New features must include unit tests.
* ✍️ **Security:** All commits must be **signed** (GPG/SSH).

## ❤️ Acknowledgments

* **[AbsaOSS](https://github.com/AbsaOSS):** For the original foundation of this project.
* **Family:** A special thanks to my wife, daughter, and son for their patience while I work in "geek mode."

## 📄 License

This project is licensed under the **Apache License 2.0**. See the [LICENSE](LICENSE) file for details.