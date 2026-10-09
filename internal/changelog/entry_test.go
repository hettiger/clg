package changelog_test

import (
	"os"
	"path/filepath"
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
			data, err := os.ReadFile(filepath.Join("testdata", tt.dataFile))
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
		name   string
		entry  changelog.ChangelogEntry
		uuidV7 func() string
		want   string
	}{
		{
			name: "changed",
			entry: changelog.ChangelogEntry{
				Type:   "changed",
				Title:  "fake message",
				Branch: "fake-branch",
			},
			uuidV7: func() string {
				return "fake-uuidv7"
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
			uuidV7: func() string {
				return "fake-uuidv7"
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
			uuidV7: func() string {
				return "fake-uuidv7"
			},
			want: "front-added-fake-uuidv7.yml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.entry.Filename(tt.uuidV7)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEntryFilenameSkipsEmptyParts(t *testing.T) {
	for _, tt := range []struct {
		group, changeType, want string
	}{
		{"ööö", "fixed", "fixed-fake-uuid.yml"},
		{"backend", "...", "backend-fake-uuid.yml"},
		{"ööö", "...", "fake-uuid.yml"},
	} {
		t.Run(tt.want, func(t *testing.T) {
			entry := changelog.ChangelogEntry{Group: tt.group, Type: tt.changeType}
			name := entry.Filename(uuidV7)
			require.Equal(t, tt.want, name)
		})
	}
}

func TestEntryYAML(t *testing.T) {
	dataValid, err := os.ReadFile(filepath.Join("testdata", "entry_valid.yml"))
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
