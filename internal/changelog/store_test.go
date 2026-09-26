package changelog_test

import (
	"testing"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/hettiger/clg/internal/testdata"
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
					Type:   "added",
					Title:  "Simple Change",
					Branch: "fake-branch",
				},
			},
		},
		{
			name:       "unexpected subdir",
			fixtureDir: "testdata/roots/unexpected_subdir",
			want: []changelog.ChangelogEntry{
				{
					Type:   "added",
					Title:  "Simple Change",
					Branch: "fake-branch",
				},
			},
		},
		{
			name:       "single entry",
			fixtureDir: "testdata/roots/single_entry",
			want: []changelog.ChangelogEntry{
				{
					Type:   "added",
					Title:  "Simple Change",
					Branch: "fake-branch",
				},
			},
		},
		{
			name:       "multiple entries",
			fixtureDir: "testdata/roots/multiple_entries",
			want: []changelog.ChangelogEntry{
				{
					Type:   "added",
					Title:  "New Feature",
					Branch: "fake-branch",
				},
				{
					Type:   "changed",
					Title:  "Changed Feature",
					Branch: "fake-branch",
				},
				{
					Type:   "fixed",
					Title:  "Fixed Feature",
					Branch: "fake-branch",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := changelog.NewEntryStore(
				tt.fixtureDir,
				uuidV7,
				map[string]string{},
				testdata.Types(),
			)

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
		groups   map[string]string
		entry    changelog.ChangelogEntry
		wantPath string
		wantErr  bool
	}{
		{
			name: "valid change",
			entry: changelog.ChangelogEntry{
				Type:   "changed",
				Title:  "Simple Change",
				Branch: "fake-branch",
			},
			wantPath: "changelogs/unreleased/changed-fake-uuid.yml",
		},
		{
			name: "valid feature",
			entry: changelog.ChangelogEntry{
				Type:   "added",
				Title:  "Simple Feature",
				Branch: "fake-branch",
			},
			wantPath: "changelogs/unreleased/added-fake-uuid.yml",
		},
		{
			name:   "valid group",
			groups: testdata.Groups(),
			entry: changelog.ChangelogEntry{
				Group:  "front",
				Type:   "added",
				Title:  "Simple Feature",
				Branch: "fake-branch",
			},
			wantPath: "changelogs/unreleased/front-added-fake-uuid.yml",
		},
		{
			name: "invalid type",
			entry: changelog.ChangelogEntry{
				Type:   "invalid",
				Title:  "Simple Feature",
				Branch: "fake-branch",
			},
			wantErr: true,
		},
		{
			name: "invalid group without config",
			entry: changelog.ChangelogEntry{
				Group:  "invalid",
				Type:   "added",
				Title:  "Simple Feature",
				Branch: "fake-branch",
			},
			wantErr: true,
		},
		{
			name:   "invalid group with config",
			groups: testdata.Groups(),
			entry: changelog.ChangelogEntry{
				Group:  "invalid",
				Type:   "added",
				Title:  "Simple Feature",
				Branch: "fake-branch",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			s := changelog.NewEntryStore(
				root,
				uuidV7,
				tt.groups,
				testdata.Types(),
			)

			gotPath, err := s.Write(tt.entry)

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, gotPath)

				return
			}

			require.NoError(t, err)
			require.Contains(t, gotPath, tt.wantPath)
		})
	}
}
