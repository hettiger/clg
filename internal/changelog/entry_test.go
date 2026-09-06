package changelog_test

import (
	"os"
	"testing"
	"time"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntryFilenameAt(t *testing.T) {
	tests := []struct {
		name  string
		entry changelog.Entry
		t     time.Time
		want  string
	}{
		{
			name:  "changed",
			entry: changelog.Entry{Title: "fake message", Type: "changed"},
			t:     time.Date(2026, 9, 5, 16, 32, 57, 0, time.UTC),
			want:  "2026-09-05-163257-changed.yml",
		},
		{
			name:  "added",
			entry: changelog.Entry{Title: "fake message", Type: "added"},
			t:     time.Date(2026, 9, 5, 16, 32, 57, 0, time.UTC),
			want:  "2026-09-05-163257-added.yml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.entry.FilenameAt(tt.t)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEntryFromYAML(t *testing.T) {
	tests := []struct {
		name     string
		dataFile string
		want     changelog.Entry
		wantErr  bool
	}{
		{
			name:     "valid",
			dataFile: "entry_valid.yml",
			want: changelog.Entry{
				Title:  "Fake Title",
				Type:   "added",
				Author: "Fake Author",
				Group:  "Fake Group",
			},
			wantErr: false,
		},
		{
			name:     "invalid",
			dataFile: "entry_invalid.yml",
			want:     changelog.Entry{},
			wantErr:  true,
		},
		{
			name:     "unsupported type",
			dataFile: "entry_unsupported_type.yml",
			want:     changelog.Entry{},
			wantErr:  true,
		},
		{
			name:     "empty title",
			dataFile: "entry_empty_title.yml",
			want:     changelog.Entry{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile("testdata/" + tt.dataFile)
			require.NoError(t, err)
			got, gotErr := changelog.EntryFromYAML(data)

			if tt.wantErr {
				require.Error(t, gotErr)
			} else {
				require.NoError(t, gotErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEntryYAML(t *testing.T) {
	dataValid, err := os.ReadFile("testdata/entry_valid.yml")
	require.NoError(t, err)
	entry, err := changelog.EntryFromYAML(dataValid)
	require.NoError(t, err)

	gotData, err := entry.YAMLData()
	require.NoError(t, err)
	assert.Equal(t, dataValid, gotData)

	gotYAML, err := entry.YAML()
	require.NoError(t, err)
	assert.Equal(t, string(dataValid), gotYAML)
}
