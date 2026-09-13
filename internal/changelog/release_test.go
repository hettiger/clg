package changelog_test

import (
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
			got, gotErr := changelog.NewRelease(tag, tt.unreleasedEntries, time.Now())

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
