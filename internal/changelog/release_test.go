package changelog_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/hettiger/clg/internal/testdata"
	"github.com/stretchr/testify/require"
)

func TestNewRelease(t *testing.T) {
	tests := []struct {
		name         string
		groups       map[string]string
		types        map[string]string
		entries      []changelog.ChangelogEntry
		issuePrefix  string
		issuePattern string
		wantErrMsg   string
	}{
		{
			name:       "empty sections",
			groups:     map[string]string{},
			types:      map[string]string{},
			entries:    []changelog.ChangelogEntry{},
			wantErrMsg: "no sections available",
		},
		{
			name:   "valid type",
			groups: map[string]string{},
			types:  testdata.Types(),
			entries: []changelog.ChangelogEntry{
				{
					Type:   "changed",
					Title:  "Simple Change",
					Branch: "fake-branch",
				},
			},
		},
		{
			name:   "invalid type",
			groups: map[string]string{},
			types:  testdata.Types(),
			entries: []changelog.ChangelogEntry{
				{
					Type:   "invalid",
					Title:  "Simple Change",
					Branch: "fake-branch",
				},
			},
			wantErrMsg: `unknown entry type: "invalid"`,
		},
		{
			name:   "empty type",
			groups: map[string]string{},
			types:  testdata.Types(),
			entries: []changelog.ChangelogEntry{
				{
					Title:  "Simple Change",
					Branch: "fake-branch",
				},
			},
			wantErrMsg: `missing entry type`,
		},
		{
			name:   "valid group",
			groups: testdata.Groups(),
			types:  testdata.Types(),
			entries: []changelog.ChangelogEntry{
				{
					Group:  "front",
					Type:   "changed",
					Title:  "Simple Change",
					Branch: "fake-branch",
				},
			},
		},
		{
			name:   "invalid group",
			groups: testdata.Groups(),
			types:  testdata.Types(),
			entries: []changelog.ChangelogEntry{
				{
					Group:  "invalid",
					Type:   "changed",
					Title:  "Simple Change",
					Branch: "fake-branch",
				},
			},
			wantErrMsg: `unknown entry group: "invalid"`,
		},
		{
			name:   "empty group",
			groups: testdata.Groups(),
			types:  testdata.Types(),
			entries: []changelog.ChangelogEntry{
				{
					Type:   "changed",
					Title:  "Simple Change",
					Branch: "fake-branch",
				},
			},
			wantErrMsg: `missing entry group`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tag := "v0.0.0"
			got, gotErr := changelog.NewRelease(
				tag,
				tt.entries,
				time.Now(),
				tt.groups,
				tt.types,
				tt.issuePrefix,
				tt.issuePattern,
			)

			if tt.wantErrMsg != "" {
				require.EqualError(t, gotErr, tt.wantErrMsg)
				require.Empty(t, got)

				return
			}

			require.NoError(t, gotErr)
			require.NotEmpty(t, got)
		})
	}
}

func TestReleaseMarkdown(t *testing.T) {
	tests := []struct {
		name         string
		tag          string
		entries      []changelog.ChangelogEntry
		groups       map[string]string
		types        map[string]string
		issuePrefix  string
		issuePattern string
		time         time.Time
		wantFixture  string
		wantErr      bool
	}{
		{
			name: "single change",
			tag:  "v0.0.0",
			entries: []changelog.ChangelogEntry{
				{
					Type:   "changed",
					Title:  "Simple Change",
					Branch: "fake-branch",
				},
			},
			types:       testdata.Types(),
			time:        time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture: "release_single_change.md",
		},
		{
			name: "multiple changes",
			tag:  "v0.1.0",
			entries: []changelog.ChangelogEntry{
				{
					Type:   "changed",
					Title:  "First Change",
					Branch: "fake-branch",
				},
				{
					Type:   "changed",
					Title:  "Second Change",
					Branch: "fake-branch",
				},
			},
			types:       testdata.Types(),
			time:        time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture: "release_multiple_changes.md",
		},
		{
			name: "type sections",
			tag:  "v0.2.0",
			entries: []changelog.ChangelogEntry{
				{
					Type:   "changed",
					Title:  "First Change",
					Branch: "fake-branch",
				},
				{
					Type:   "fixed",
					Title:  "Simple Bug Fix",
					Branch: "fake-branch",
				},
				{
					Type:   "changed",
					Title:  "Second Change",
					Branch: "fake-branch",
				},
			},
			types:       testdata.Types(),
			time:        time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture: "release_type_sections.md",
		},
		{
			name:        "empty",
			tag:         "v0.3.1",
			entries:     []changelog.ChangelogEntry{},
			types:       testdata.Types(),
			time:        time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture: "release_empty.md",
		},
		{
			name: "groups",
			tag:  "v0.4.0",
			entries: []changelog.ChangelogEntry{
				{
					Group:  "front",
					Type:   "changed",
					Title:  "First Change",
					Branch: "fake-branch",
				},
				{
					Group:  "back",
					Type:   "fixed",
					Title:  "Important Bug Fix",
					Branch: "fake-branch",
				},
				{
					Group:  "front",
					Type:   "fixed",
					Title:  "Simple Bug Fix",
					Branch: "fake-branch",
				},
				{
					Group:  "front",
					Type:   "changed",
					Title:  "Second Change",
					Branch: "fake-branch",
				},
			},
			groups:      testdata.Groups(),
			types:       testdata.Types(),
			time:        time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture: "release_groups.md",
		},

		{
			name: "groups single change",
			tag:  "v0.5.0",
			entries: []changelog.ChangelogEntry{
				{
					Group:  "front",
					Type:   "changed",
					Title:  "Single Change",
					Branch: "fake-branch",
				},
			},
			groups:      testdata.Groups(),
			types:       testdata.Types(),
			time:        time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture: "release_groups_single_change.md",
		},

		{
			name: "single change author",
			tag:  "v0.5.1",
			entries: []changelog.ChangelogEntry{
				{
					Type:   "changed",
					Title:  "Single Change",
					Branch: "fake-branch",
					Author: "Fake Author",
				},
			},
			types:       testdata.Types(),
			time:        time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture: "release_single_change_author.md",
		},

		{
			name: "groups single change author",
			tag:  "v0.5.2",
			entries: []changelog.ChangelogEntry{
				{
					Group:  "front",
					Type:   "changed",
					Title:  "Single Change",
					Branch: "fake-branch",
					Author: "[Fake Author](https://github.com/fake-author)",
				},
			},
			groups:      testdata.Groups(),
			types:       testdata.Types(),
			time:        time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture: "release_groups_single_change_author.md",
		},

		{
			name: "single change issue",
			tag:  "v0.5.3",
			entries: []changelog.ChangelogEntry{
				{
					Type:   "changed",
					Title:  "Single Change",
					Branch: "feature/3-fake-branch",
				},
			},
			types:        testdata.Types(),
			issuePrefix:  "#",
			issuePattern: "^feature/(\\d+)-",
			time:         time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture:  "release_single_change_issue.md",
		},

		{
			name: "groups single change issue author",
			tag:  "v0.5.4",
			entries: []changelog.ChangelogEntry{
				{
					Group:  "front",
					Type:   "changed",
					Title:  "Single Change",
					Branch: "feature/3-fake-branch",
					Author: "[Fake Author](https://github.com/fake-author)",
				},
			},
			groups:       testdata.Groups(),
			types:        testdata.Types(),
			issuePrefix:  "#",
			issuePattern: "^feature/(\\d+)-",
			time:         time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture:  "release_groups_single_change_issue_author.md",
		},

		{
			name: "issue pattern mismatch",
			tag:  "v0.5.5",
			entries: []changelog.ChangelogEntry{
				{
					Type:   "changed",
					Title:  "Single Change",
					Branch: "3-fake-branch",
					Author: "Fake Author",
				},
			},
			types:        testdata.Types(),
			issuePrefix:  "#",
			issuePattern: "^feature/(\\d+)-",
			time:         time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture:  "release_issue_pattern_mismatch.md",
		},

		{
			name: "issue pattern invalid",
			tag:  "v0.5.6",
			entries: []changelog.ChangelogEntry{
				{
					Type:   "changed",
					Title:  "Single Change",
					Branch: "feature/3-fake-branch",
				},
			},
			types:        testdata.Types(),
			issuePrefix:  "#",
			issuePattern: "^feature/(\\d+-",
			time:         time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantErr:      true,
		},

		{
			name: "invalid capturing groups",
			tag:  "v0.5.7",
			entries: []changelog.ChangelogEntry{
				{
					Type:   "changed",
					Title:  "Single Change",
					Branch: "feature/3-fake-branch",
				},
			},
			types:        testdata.Types(),
			issuePrefix:  "#",
			issuePattern: "^feature/(\\d+)-(.+)$",
			time:         time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			release, err := changelog.NewRelease(
				tt.tag,
				tt.entries,
				tt.time,
				tt.groups,
				tt.types,
				tt.issuePrefix,
				tt.issuePattern,
			)
			require.NoError(t, err)

			var want string
			if !tt.wantErr {
				wantData, err := os.ReadFile(filepath.Join("testdata", tt.wantFixture))
				require.NoError(t, err)
				want = strings.TrimSpace(string(wantData))
			}

			got, gotErr := release.Markdown()

			if tt.wantErr {
				require.Error(t, gotErr)
			} else {
				require.NoError(t, gotErr)
			}

			require.Equal(t, want, got)
		})
	}
}
