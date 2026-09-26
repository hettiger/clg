package changelog_test

import (
	"errors"
	"os"
	"testing"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/hettiger/clg/internal/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewChangelogEntry(t *testing.T) {
	tests := []struct {
		name      string
		dataFile  string
		groupKeys []string
		want      changelog.ChangelogEntry
		wantErr   bool
	}{
		{
			name:     "valid",
			dataFile: "entry_valid.yml",
			want: changelog.ChangelogEntry{
				Group:  "",
				Type:   "added",
				Title:  "Fake Title",
				Author: "Fake Author",
				Branch: "fake-branch",
			},
			wantErr: false,
		},
		{
			name:      "valid with group",
			dataFile:  "entry_valid_group.yml",
			groupKeys: testdata.GroupKeys(),
			want: changelog.ChangelogEntry{
				Group:  "front",
				Type:   "added",
				Title:  "Fake Title",
				Author: "Fake Author",
				Branch: "fake-branch",
			},
			wantErr: false,
		},
		{
			name:     "invalid",
			dataFile: "entry_invalid.yml",
			want:     changelog.ChangelogEntry{},
			wantErr:  true,
		},
		{
			name:     "unsupported group without groups configured",
			dataFile: "entry_unsupported_group.yml",
			want:     changelog.ChangelogEntry{},
			wantErr:  true,
		},
		{
			name:      "unsupported group with groups configured",
			groupKeys: testdata.GroupKeys(),
			dataFile:  "entry_unsupported_group.yml",
			want:      changelog.ChangelogEntry{},
			wantErr:   true,
		},
		{
			name:     "unsupported type",
			dataFile: "entry_unsupported_type.yml",
			want:     changelog.ChangelogEntry{},
			wantErr:  true,
		},
		{
			name:     "empty title",
			dataFile: "entry_empty_title.yml",
			want:     changelog.ChangelogEntry{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile("testdata/" + tt.dataFile)
			require.NoError(t, err)
			got, gotErr := changelog.NewChangelogEntry(data, tt.groupKeys, testdata.TypeKeys())

			if tt.wantErr {
				require.Error(t, gotErr)
			} else {
				require.NoError(t, gotErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEntryFilename(t *testing.T) {
	tests := []struct {
		name    string
		entry   changelog.ChangelogEntry
		uuidV7  func() (string, error)
		want    string
		wantErr bool
	}{
		{
			name: "changed",
			entry: changelog.ChangelogEntry{
				Type:   "changed",
				Title:  "fake message",
				Branch: "fake-branch",
			},
			uuidV7: func() (string, error) {
				return "fake-uuidv7", nil
			},
			want: "changed-fake-uuidv7.yml",
		},
		{
			name: "added",
			entry: changelog.ChangelogEntry{
				Type:   "added",
				Title:  "fake message",
				Branch: "fake-branch",
			},
			uuidV7: func() (string, error) {
				return "fake-uuidv7", nil
			},
			want: "added-fake-uuidv7.yml",
		},
		{
			name: "front fixed",
			entry: changelog.ChangelogEntry{
				Group:  "front",
				Type:   "added",
				Title:  "fake message",
				Branch: "fake-branch",
			},
			uuidV7: func() (string, error) {
				return "fake-uuidv7", nil
			},
			want: "front-added-fake-uuidv7.yml",
		},
		{
			name: "uuid error",
			entry: changelog.ChangelogEntry{
				Group:  "front",
				Type:   "added",
				Title:  "fake message",
				Branch: "fake-branch",
			},
			uuidV7: func() (string, error) {
				return "", errors.New("error fake")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := tt.entry.Filename(tt.uuidV7)

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
	entry, err := changelog.NewChangelogEntry(dataValid, []string{}, testdata.TypeKeys())
	require.NoError(t, err)

	gotData, err := entry.YAMLData()
	require.NoError(t, err)
	assert.Equal(t, dataValid, gotData)

	gotYAML, err := entry.YAML()
	require.NoError(t, err)
	assert.Equal(t, string(dataValid), gotYAML)
}
