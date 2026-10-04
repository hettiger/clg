# clg

`clg` turns small YAML release notes into a Markdown `CHANGELOG.md`. Record
changes during development, then publish them together as a dated release.

> Pre-alpha. No releases yet; interfaces may change without notice.

## Installation

Requires Go 1.27.1 or newer.

```sh
go install github.com/hettiger/clg@latest
```

Or build from a checkout with `go build -o clg .`.

## Quick start

Add `<!-- CLG -->` to your `CHANGELOG.md`, then run from your Git working tree:

```sh
clg new -t added -m "Support exporting reports" -a "Jane Doe"
clg show
clg release v1.2.0
```

This inserts a section after the marker and removes the released entry files
from `changelogs/unreleased/`:

```md
## [v1.2.0] - 2026-09-06

### New Feature (1 change)

- Support exporting reports (Jane Doe)
```

## Commands

### `clg new`

Create an entry. Prompts for missing type, message, and group (when configured).
Supply all three to run non-interactively. Requires a Git working tree to record
the current branch.

| Flag | Description |
| --- | --- |
| `-t, --type` | Change type key. |
| `-m, --message` | Entry text. |
| `-g, --group` | Group key; available when groups are configured. |
| `-a, --author` | Author string; see [Author attribution](#author-attribution). |

### `clg show`

List unreleased entries with their type, title, author, branch, and group (if
configured). Author strings are shown literally, including Markdown syntax.
Use `-b, --branch` to filter by branch.

### `clg release <tag>`

Group entries by type, or by group and then type, and insert a release into
`CHANGELOG.md` using the current date in the configured timezone (UTC by default).
Deletes the source entries afterward; does nothing when there are no entries.

The insertion marker must already exist. Override the configured marker with
`-m, --marker` or the timezone with `-t, --timezone`.

### `clg clean`

Delete all unreleased entries after confirmation. Use `--force` to skip the
prompt. Does not modify `CHANGELOG.md`.

### `clg help`

Run `clg help <command>` or `clg <command> --help` for command-line help.

## Configuration

Both configuration files are optional:

- `~/.clg.yml`: personal settings across projects.
- `.clg.yml` in the current working directory: shared project settings.

Project settings override global settings, which override built-in defaults.
Nested mappings merge field by field. Environment variables override file
settings, using the `CLG_` prefix and underscores for nested keys (for example,
`CLG_MARKER` and `CLG_TIMEZONE`). Empty environment values are ignored, except
for `CLG_AUTHOR`.

Set `timezone` to an IANA location name such as `America/New_York`; it defaults
to `UTC`.

The default marker, timezone, and change types are:

```yaml
marker: "<!-- CLG -->"
timezone: UTC
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

Use `types` to customize keys and headings. Add `groups` to require a group for
each entry and render group headings above type headings:

```yaml
groups:
  front: Frontend
  back: Backend
```

### Author attribution

Authors are strings: plain names, `@mentions`, or Markdown links. Keep your
identity in `~/.clg.yml` rather than the shared project configuration:

```yaml
author: "[Jane Doe](https://example.com/jane)"
```

Precedence, highest first: `--author` → `CLG_AUTHOR` → project `author` → global
`author` → Git's `user.name`.

An explicit empty string overrides lower-priority sources: use `author: ""`
for a default, `CLG_AUTHOR=""` for an environment override, or `--author ""`
for one entry. Git lookup is best-effort; missing authors are silently omitted.
`clg new` trims surrounding whitespace and never prompts for an author.

Releases append the recorded author in parentheses, verbatim. Markdown is not
escaped and links are not validated. Attribution comes from the entry, not the
configuration at release time.

## Entry format

Entries live in `changelogs/unreleased/`, named `<type>-<UUIDv7>.yml` or
`<group>-<type>-<UUIDv7>.yml` when groups are configured:

```yaml
group: back
type: changed
title: Improve report permissions
author: "[Jane Doe](https://example.com/jane)"
branch: feature/report-export
```

Every entry needs a non-empty `title`, a configured `type`, and a configured
`group` when groups are enabled. `author` is optional; `clg new` omits it when
empty. Invalid files block commands that read unreleased entries.

## Development

```sh
go test ./...
```

## License

See [LICENSE](LICENSE).
