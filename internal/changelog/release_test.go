package changelog_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestNewRelease(t *testing.T) {
	tests := []struct {
		name              string
		unreleasedEntries []changelog.ChangelogEntry
		wantErrMsg        string
	}{
		{
			name: "valid type",
			unreleasedEntries: []changelog.ChangelogEntry{
				{
					Title: "Simple Change",
					Type:  "changed",
				},
			},
		},
		{
			name: "invalid type",
			unreleasedEntries: []changelog.ChangelogEntry{
				{
					Title: "Simple Change",
					Type:  "invalid",
				},
			},
			wantErrMsg: "Unknown type keyword provided (invalid)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tag := "v0.0.0"
			got, gotErr := changelog.NewRelease(tag, tt.unreleasedEntries, time.Now(), testTypes())

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
		name              string
		tag               string
		unreleasedEntries []changelog.ChangelogEntry
		time              time.Time
		wantFixture       string
	}{
		{
			name: "single change",
			tag:  "v0.0.0",
			unreleasedEntries: []changelog.ChangelogEntry{
				{
					Title: "Simple Change",
					Type:  "changed",
				},
			},
			time:        time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture: "release_single_change.md",
		},
		{
			name: "multiple changes",
			tag:  "v0.1.0",
			unreleasedEntries: []changelog.ChangelogEntry{
				{
					Title: "First Change",
					Type:  "changed",
				},
				{
					Title: "Second Change",
					Type:  "changed",
				},
			},
			time:        time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture: "release_multiple_changes.md",
		},
		{
			name: "groups",
			tag:  "v0.2.0",
			unreleasedEntries: []changelog.ChangelogEntry{
				{
					Title: "First Change",
					Type:  "changed",
				},
				{
					Title: "Simple Bug Fix",
					Type:  "fixed",
				},
				{
					Title: "Second Change",
					Type:  "changed",
				},
			},
			time:        time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture: "release_groups.md",
		},
		{
			name: "ignored",
			tag:  "v0.3.0",
			unreleasedEntries: []changelog.ChangelogEntry{
				{
					Title: "First Change",
					Type:  "changed",
				},
				{
					Title: "Simple Bug Fix",
					Type:  "fixed",
				},
				{
					Title: "Ignored Change",
					Type:  "ignore",
				},
				{
					Title: "Second Change",
					Type:  "changed",
				},
			},
			time:        time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture: "release_ignored.md",
		},
		{
			name:              "empty",
			tag:               "v0.3.1",
			unreleasedEntries: []changelog.ChangelogEntry{},
			time:              time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC),
			wantFixture:       "release_empty.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			release, err := changelog.NewRelease(tt.tag, tt.unreleasedEntries, tt.time, testTypes())
			require.NoError(t, err)
			wantData, err := os.ReadFile(filepath.Join("testdata", tt.wantFixture))
			require.NoError(t, err)
			want := strings.TrimSpace(string(wantData))

			got := release.Markdown()

			require.Equal(t, want, got)
		})
	}
}
