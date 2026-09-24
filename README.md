# clg

`clg` is a small command-line tool for maintaining a Markdown `CHANGELOG.md`.
Instead of editing the changelog during development, record each change as a
YAML file and turn all unreleased entries into a dated release when you publish.

> **Warning:** This project is work in progress. Use at your own risk. APIs might change.
> The docs are 100% AI-generated. No releases published yet.

## How it works

`clg` uses this layout in the current working directory:

```text
.
├── CHANGELOG.md
├── clg.yml                  # optional configuration
└── changelogs/
    └── unreleased/
        └── added-0199321f-7b2c-7c4f-bd12-4c5f8f7c2a10.yml
```

The files in `changelogs/unreleased/` are temporary release notes. `clg release`
groups them by type—or by group and then type when groups are configured—inserts
the resulting Markdown into `CHANGELOG.md`, and removes the source files.

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

Add the insertion marker to `CHANGELOG.md` once, usually near the top. The
default marker is `<!-- CLG -->`:

```md
<!-- CLG -->
```

Optionally configure groups, custom labels, or a different marker in `clg.yml`:

```yaml
marker: "<!-- CLG -->"
groups:
  front: Frontend
  back: Backend
types:
  added: New Feature
  fixed: Bug Fix
```

When `groups` is configured, every entry must specify one of its keys and
releases are rendered with group headings containing type headings.

Record a change. With no flags, `clg new` asks for the configured group (if any),
type, and message:

```sh
clg new
```

For scripts or a faster workflow, provide both values directly:

```sh
clg new --type added --message "Support exporting reports"
clg new -g back -t fixed -m "Prevent duplicate notifications"
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

With groups configured, the release uses one additional heading level:

```md
### Backend

#### Bug Fix (1 change)

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
| `-g, --group` | Configured group key. If omitted, choose from an interactive list when groups are configured. |
| `-t, --type` | Configured change type. If omitted, choose from an interactive list. |
| `-m, --message` | Entry text. If omitted, enter it interactively. |

The flags can be supplied together, which makes the command non-interactive.
The generated filename contains the type and a UUIDv7, for example
`fixed-0199321f-7b2c-7c4f-bd12-4c5f8f7c2a10.yml`. When groups are configured,
the group key is prefixed to the filename, for example
`back-fixed-0199321f-7b2c-7c4f-bd12-4c5f8f7c2a10.yml`.

### `clg show`

Display all valid entries that have not yet been released:

```sh
clg show
```

The output includes the type, message, and optional author. Group values are
stored in entries and used when generating a release. If there are no entries,
`clg` reports that there is nothing to show.

### `clg release [tag]`

Convert all unreleased entries into a release and insert it into
`CHANGELOG.md`:

```sh
clg release v1.2.0
```

| Flag | Default | Description |
| --- | --- | --- |
| `-m, --marker` | configured marker or `<!-- CLG -->` | Text where the new release is inserted. |

The marker must already exist in `CHANGELOG.md`. To use a different marker:

```sh
clg release v1.2.0 --marker "<!-- RELEASES -->"
```

If there are no unreleased entries, the command leaves the changelog unchanged.

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

The heading is used when `clg release` groups entries.

## Configuration

Configuration is loaded from an optional `clg.yml` in the current working
directory. The defaults are:

```yaml
marker: "<!-- CLG -->"
types:
  added: New Feature
  fixed: Bug Fix
  hotfix: Hotfix
  changed: Feature Change
  deprecated: New Deprecation
  removed: Feature Removal
  security: Security Fix
  performance: Performance Improvement
  other: Other
```

Use `groups` to enable grouped releases. Group keys are used in entry files and
CLI flags; their values are the headings shown in generated Markdown:

```yaml
groups:
  front: Frontend
  back: Backend
```

## Entry format

Each entry is a YAML document. `clg new` writes the required `title` and `type`
fields, plus `group` when groups are configured:

```yaml
title: Improve report permissions
type: changed
group: back
author: Jane Doe
```

The `author` field is optional and is displayed by `clg show`; it is not set by
`clg new`. Every YAML file in `changelogs/unreleased/` must have a non-empty
`title`, a configured `type`, and—when groups are configured—a configured
`group`. Invalid files prevent commands that read unreleased entries from
completing.

## Typical release workflow

```sh
# During development
clg new -g back -t added -m "Add CSV export"
clg new -g front -t fixed -m "Handle empty report filters"

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
