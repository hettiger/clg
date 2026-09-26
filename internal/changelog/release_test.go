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
		name       string
		groups     map[string]string
		types      map[string]string
		entries    []changelog.ChangelogEntry
		wantErrMsg string
	}{
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
					Type:   "",
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
					Group:  "",
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
		name        string
		tag         string
		entries     []changelog.ChangelogEntry
		groups      map[string]string
		types       map[string]string
		time        time.Time
		wantFixture string
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
			)
			require.NoError(t, err)
			wantData, err := os.ReadFile(filepath.Join("testdata", tt.wantFixture))
			require.NoError(t, err)
			want := strings.TrimSpace(string(wantData))

			got := release.Markdown()

			require.Equal(t, want, got)
		})
	}
}
