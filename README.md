# ghostchar

CLI tool that detects invisible Unicode characters, Private Use Area (PUA) codepoints, and bidirectional control characters in source code files.

Built for local development workflows, pre-commit hooks, and CI/CD pipelines. Ships as a single static binary with zero dependencies.

## Why

Invisible Unicode characters can hide in source code and cause subtle bugs, broken builds, or security vulnerabilities. Bidirectional control characters (U+202A–U+202E, U+2066–U+2069) are a known [supply-chain attack vector](https://trojansource.codes/) that can make code appear different from what compilers actually execute. `ghostchar` catches these before they reach production.

## Installation

### Homebrew

```bash
brew install sgaunet/tap/ghostchar
```

### Docker

```bash
docker run --rm -v "$PWD":/src ghcr.io/sgaunet/ghostchar:latest scan /src
```

### Binary

Download a prebuilt binary from [GitHub Releases](https://github.com/sgaunet/ghostchar/releases) for your platform:

- `linux/amd64`, `linux/arm64`
- `darwin/amd64`, `darwin/arm64`
- `windows/amd64`

## Quick Start

```bash
# Scan current directory
ghostchar .

# Scan specific files
ghostchar scan main.go utils.go

# Scan only Go and Python files
ghostchar scan --ext go,py ./src

# JSON output for CI integration
ghostchar scan --format json .

# SARIF output for GitHub Advanced Security / GitLab SAST
ghostchar scan --format sarif . > results.sarif

# Only detect bidi control characters
ghostchar scan --categories bidi .

# Quiet mode — exit code only
ghostchar scan -q .
```

## Usage

```
ghostchar [command] [flags] [path...]
```

### Commands

| Command | Description |
|---------|-------------|
| `scan` | Scan files or directories (default command) |
| `list-chars` | Print all detected character categories and codepoints |
| `version` | Print version, build date, and commit hash |

### Flags (`scan`)

| Flag | Default | Description |
|------|---------|-------------|
| `--ext` | common source extensions | Comma-separated file extensions to scan (e.g. `go,py,js`) |
| `--exclude` | `.git,vendor,node_modules` | Comma-separated directories to exclude |
| `--categories` | `all` | Categories to detect: `invisible`, `pua`, `bidi`, `all` |
| `--format` | `text` | Output format: `text`, `json`, `sarif` |
| `--quiet` / `-q` | `false` | Suppress output, only use exit code |
| `--no-color` | `false` | Disable ANSI color output |
| `--max-file-size` | `1MB` | Skip files larger than this size |

### Default Scanned Extensions

```
go, py, js, ts, jsx, tsx, java, c, cpp, h, hpp, cs, rb, php,
rs, kt, swift, sh, bash, zsh, yaml, yml, toml, json, xml, html,
htm, css, scss, sql, tf, md, txt
```

## Detected Characters

### Invisible Characters

| Codepoint | Name |
|-----------|------|
| U+200B | Zero-Width Space |
| U+200C | Zero-Width Non-Joiner |
| U+200D | Zero-Width Joiner |
| U+2060 | Word Joiner |
| U+2061 | Function Application |
| U+2062 | Invisible Times |
| U+2063 | Invisible Separator |
| U+2064 | Invisible Plus |
| U+FEFF | BOM / Zero-Width No-Break Space |
| U+00AD | Soft Hyphen |
| U+034F | Combining Grapheme Joiner |

### Private Use Area (PUA)

| Range | Block |
|-------|-------|
| U+E000 – U+F8FF | BMP Private Use Area |
| U+F0000 – U+FFFFF | Supplementary PUA-A |
| U+100000 – U+10FFFF | Supplementary PUA-B |

### Bidirectional Control Characters

| Codepoint | Name |
|-----------|------|
| U+202A–U+202E | LTR/RTL overrides and embeddings |
| U+2066–U+2069 | Isolate and PDF controls |
| U+200E–U+200F | LTR / RTL marks |

Use `ghostchar list-chars` to print the full table at any time.

## Output Formats

### Text (default)

```
path/to/file.go:12:5  U+200B  ZERO WIDTH SPACE         [invisible]
path/to/file.go:34:1  U+202E  RIGHT-TO-LEFT OVERRIDE   [bidi]

2 findings in 1 file (47 files scanned)
```

### JSON

```json
{
  "summary": {
    "files_scanned": 47,
    "files_with_findings": 1,
    "total_findings": 2
  },
  "findings": [
    {
      "file": "path/to/file.go",
      "line": 12,
      "column": 5,
      "codepoint": "U+200B",
      "name": "ZERO WIDTH SPACE",
      "category": "invisible"
    }
  ]
}
```

### SARIF

SARIF 2.1.0 format for integration with GitHub Advanced Security, GitLab SAST, and compatible tools.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | No findings |
| 1 | One or more findings detected |
| 2 | Runtime error (I/O error, invalid flags, etc.) |

## CI/CD Integration

### GitHub Actions

```yaml
- name: Check for invisible characters
  run: ghostchar scan --format sarif . > ghostchar.sarif

- name: Upload SARIF
  uses: github/codeql-action/upload-sarif@v3
  with:
    sarif_file: ghostchar.sarif
```

### Pre-commit Hook

```bash
#!/bin/sh
ghostchar scan --quiet .
```

### GitLab CI

```yaml
ghostchar:
  image: ghcr.io/sgaunet/ghostchar:latest
  script:
    - ghostchar scan --format json .
```

## File Processing

- Files are read as UTF-8; invalid byte sequences are skipped with a warning
- Binary files are detected heuristically (null bytes in first 512 bytes) and skipped
- Symlinks are not followed

## Project Structure

```
.
├── cmd/
│   └── root.go         # cobra root + scan command
├── internal/
│   ├── scanner/        # file walking + character detection logic
│   ├── report/         # output formatters (text, json, sarif)
│   └── charset/        # codepoint definitions and category logic
├── main.go
├── .goreleaser.yaml
├── Dockerfile
└── README.md
```

## Roadmap

- `--fix` flag to strip or replace detected characters in-place
- UTF-16 / Latin-1 support with automatic encoding detection
- Homoglyph / confusable character detection ([UAX #39](https://unicode.org/reports/tr39/))
- `.ghostignore` config file for per-repo allowlists

## License

[MIT](LICENSE)
