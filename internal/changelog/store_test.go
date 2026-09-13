package changelog_test

import (
	"testing"
	"time"

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

	now := func() time.Time {
		return time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC)
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
