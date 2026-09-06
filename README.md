# clg

`clg` is a small command-line tool for maintaining a Markdown `CHANGELOG.md`.
Instead of editing the changelog during development, record each change as a
YAML file and turn all unreleased entries into a dated release when you publish.

> **Warning:** This project is work in progress. Use at your own risk. APIs might change.
> Not all features are fully implemented yet (authors, groups) but the app is already usable.
> The docs are 100% AI-generated. No releases published yet.

## How it works

`clg` uses this layout in the current working directory:

```text
.
├── CHANGELOG.md
└── changelogs/
    └── unreleased/
        └── 2026-09-06-142530-added.yml
```

The files in `changelogs/unreleased/` are temporary release notes. `clg release`
groups them by type, inserts the resulting Markdown into `CHANGELOG.md`, and
removes the source files.

## Installation

With Go installed:

```sh
go install github.com/hettiger/clg@latest
```

Or build the binary from a checkout:

```sh
go build -o clg .
```

`clg` requires Go 1.25.8 or newer.

## Quick start

Add the insertion marker to `CHANGELOG.md` once, usually near the top:

```md
<!-- CLG -->
```

Record a change. With no flags, `clg new` asks for the type and message:

```sh
clg new
```

For scripts or a faster workflow, provide both values directly:

```sh
clg new --type added --message "Support exporting reports"
clg new -t fixed -m "Prevent duplicate notifications"
```

Review the unreleased entries:

```sh
clg show
```

When you are ready to publish, pass the release tag:

```sh
clg release v1.2.0
```

This adds a section like the following immediately after `<!-- CLG -->`:

```md
## [v1.2.0] - 2026-09-06

### New Feature (1 change)

- Support exporting reports

### Bug Fix (1 change)

- Prevent duplicate notifications
```

The release date is the current UTC date.

## Commands

### `clg new`

Create an unreleased changelog entry in `changelogs/unreleased/`.

```sh
clg new [flags]
```

| Flag | Description |
| --- | --- |
| `-t, --type` | Change type. If omitted, choose from an interactive list. |
| `-m, --message` | Entry text. If omitted, enter it interactively. |

Both flags can be supplied together, which makes the command non-interactive.
The generated filename contains the UTC timestamp and type, for example
`2026-09-06-142530-fixed.yml`.

### `clg show`

Display all valid entries that have not yet been released:

```sh
clg show
```

The output includes the type, message, and optional author. If there are no
entries, `clg` reports that there is nothing to show.

### `clg release [tag]`

Convert all unreleased entries into a release and insert it into
`CHANGELOG.md`:

```sh
clg release v1.2.0
```

| Flag | Default | Description |
| --- | --- | --- |
| `-m, --marker` | `<!-- CLG -->` | Text where the new release is inserted. |

The marker must already exist in `CHANGELOG.md`. To use a different marker:

```sh
clg release v1.2.0 --marker "<!-- RELEASES -->"
```

If there are no unreleased entries, the command leaves the changelog
unchanged. Entries with the `ignore` type are accepted as temporary notes but
are not included in the generated release.

### `clg clean`

Delete all unreleased entry files:

```sh
clg clean
```

The command asks for confirmation. Use `--force` when confirmation is not
possible or desired:

```sh
clg clean --force
```

This only removes files in `changelogs/unreleased/`; it does not modify
`CHANGELOG.md`.

## Change types

The following type keywords are supported:

| Keyword | Heading |
| --- | --- |
| `added` | New Feature |
| `fixed` | Bug Fix |
| `hotfix` | Hotfix |
| `changed` | Feature Change |
| `deprecated` | New Deprecation |
| `removed` | Feature Removal |
| `security` | Security Fix |
| `performance` | Performance Improvement |
| `other` | Other |
| `ignore` | No Changelog |

The heading is used when `clg release` groups entries. The `ignore` type is
deliberately omitted from release Markdown.

## Entry format

Each entry is a YAML document. `clg new` writes the required `title` and
`type` fields:

```yaml
title: Support exporting reports
type: added
```

Optional `author` and `group` fields are supported when entries are created or
edited manually:

```yaml
title: Improve report permissions
type: changed
author: Jane Doe
group: Reporting
```

Every YAML file in `changelogs/unreleased/` must have a non-empty `title` and a
supported `type`; invalid files prevent commands that read unreleased entries
from completing.

## Typical release workflow

```sh
# During development
clg new -t added -m "Add CSV export"
clg new -t fixed -m "Handle empty report filters"

# Before publishing
clg show
clg release v1.2.0
git diff -- CHANGELOG.md
git add CHANGELOG.md
git commit -m "Release v1.2.0"
```

Run `clg --help` or `clg <command> --help` for the command-line help.

## Development

Run the test suite with:

```sh
go test ./...
```

## License

See [LICENSE](LICENSE).
