# 🚀 Golic

> **A declarative tool for injecting and managing license headers in source code.**

Golic automates the tedious task of ensuring every source file in your project has the correct license header. It’s built for developers who want a "set it and forget it" solution for compliance.

```bash
# Preview changes before applying
golic inject -c="2026 MyCompany ltd." --dry
```

## 📥 Installation

Install the binary using Go 1.16+:

```bash
go install github.com/BUgs5382/golic@latest
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

Contains license text and formatting rules. Golic merges your local file with its [embedded master configuration](.golic.yaml) by default.

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
      suffix: "/*}}"
    .mzm:
      prefix: "" # No indent/prefix; place text directly at top
```

## 🛠️ Usage

### 💉 Injecting Licenses

Once your config is ready, run:

```bash
golic inject -t apacheX
```

> [\!TIP]
> Always use the `--dry` flag first to see a preview of which files will be modified without actually changing them.

### 🔄 Updating or Removing

To update a license, you must remove the old one first:

1.  **Remove:** `golic remove -t apacheX`
2.  **Update:** Modify `.golic.yaml`
3.  **Re-inject:** `golic inject -t apacheX`

### 🤖 CI/CD Integration

To fail a build if licenses are missing (e.g., a developer forgot to run Golic), use the `-x` (exit code) flag:

```bash
golic inject --dry -x -t apache2
```

## 🏗 Development

### 🛠 Build

To compile the project locally, execute:

```bash
make build
```

To remove build artifacts and clean your workspace:

```bash
make clean
```

If you are **contributing** to this project, you must first initialize the linting environment:
 ```bash
 make lint-init
 ```
This command installs all necessary dependencies and tools for code analysis.

Once initialized, you can analyze the codebase by running:

```bash
make lint
```

To verify only the project licenses, use:

```bash
make license
```

### 🧪 Test

To execute the unit testing suite, run:

```bash
make test
```

## 📋 Command Reference

| Command | Description |
| :--- | :--- |
| `inject` | Injects license headers based on templates. |
| `remove` | Removes license headers matching the config. |
| `version` | Prints current version. |

**Common Flags:**

* `-c, --copyright` : Set the holder/year (Default: `YYYY [Insert Company]`). YYYY defaults to the current year.
* `-d, --dry` : Run without writing to files.
* `-t, --template` : Specify which license key to use from the config.
* `-v, --verbose` : Enable detailed trace logging.

## 🤝 Contributing

We welcome Pull Requests! Please follow these steps:

* **✅ Validation:** Run `make lint` to verify code quality.
* **🧪 Testing:** New features must include unit tests.
* **✍️ Security:** All commits must be **signed** (GPG/SSH).

## ❤️ Acknowledgments

* **[AbsaOSS](https://github.com/AbsaOSS):** For the original foundation of this project.
* **Family:** A special thanks to my wife, daughter, and son for their patience while I work in "geek mode."

## 📄 License

This project is licensed under the **Apache License 2.0**. See the [LICENSE](LICENSE) file for details.