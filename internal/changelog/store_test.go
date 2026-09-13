package changelog_test

import (
	"testing"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestEntryStoreUnreleasedEntries(t *testing.T) {
	tests := []struct {
		name       string
		fixtureDir string
		want       []changelog.ChangelogEntry
		wantErrMsg string
	}{
		{
			name:       "empty root",
			fixtureDir: t.TempDir(),
			want:       []changelog.ChangelogEntry{},
		},
		{
			name:       "invalid entry",
			fixtureDir: "testdata/roots/invalid_entry",
			wantErrMsg: "type must be one of the following values",
		},
		{
			name:       "unexpected file",
			fixtureDir: "testdata/roots/unexpected_file",
			want: []changelog.ChangelogEntry{
				{
					Title: "Simple Change",
					Type:  "added",
				},
			},
		},
		{
			name:       "unexpected subdir",
			fixtureDir: "testdata/roots/unexpected_subdir",
			want: []changelog.ChangelogEntry{
				{
					Title: "Simple Change",
					Type:  "added",
				},
			},
		},
		{
			name:       "single entry",
			fixtureDir: "testdata/roots/single_entry",
			want: []changelog.ChangelogEntry{
				{
					Title: "Simple Change",
					Type:  "added",
				},
			},
		},
		{
			name:       "multiple entries",
			fixtureDir: "testdata/roots/multiple_entries",
			want: []changelog.ChangelogEntry{
				{
					Title: "New Feature",
					Type:  "added",
				},
				{
					Title: "Changed Feature",
					Type:  "changed",
				},
				{
					Title: "Fixed Feature",
					Type:  "fixed",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := changelog.NewEntryStore(tt.fixtureDir, now)

			got, gotErr := s.UnreleasedEntries()

			if tt.wantErrMsg != "" {
				require.ErrorContains(t, gotErr, tt.wantErrMsg)
				require.Empty(t, got)

				return
			}

			require.NoError(t, gotErr)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestEntryStoreWrite(t *testing.T) {
	tests := []struct {
		name     string
		entry    changelog.ChangelogEntry
		wantPath string
		wantErr  bool
	}{
		{
			name: "valid change",
			entry: changelog.ChangelogEntry{
				Title: "Simple Change",
				Type:  "changed",
			},
			wantPath: "changelogs/unreleased/2026-09-13-182541-changed.yml",
		},
		{
			name: "valid feature",
			entry: changelog.ChangelogEntry{
				Title: "Simple Feature",
				Type:  "added",
			},
			wantPath: "changelogs/unreleased/2026-09-13-182541-added.yml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			s := changelog.NewEntryStore(root, now)

			gotPath, err := s.Write(tt.entry)

			require.NoError(t, err)
			require.Contains(t, gotPath, tt.wantPath)
		})
	}
}
