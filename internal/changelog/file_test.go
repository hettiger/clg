package changelog_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewChangelogFile(t *testing.T) {
	tests := []struct {
		name   string
		dir    string
		marker string
		want   changelog.ChangelogFile
	}{
		{
			name:   "valid input",
			dir:    "/tmp/test/",
			marker: "<!-- Fake Marker -->",
			want: changelog.ChangelogFile{
				Path:   "/tmp/test/CHANGELOG.md",
				Marker: "<!-- Fake Marker -->",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := changelog.NewChangelogFile(tt.dir, tt.marker)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestChangelogFileWriteAndRead(t *testing.T) {
	t.Parallel()

	want := "Fake Log"
	dir := t.TempDir()
	marker := "<!-- Fake Marker -->"
	file := changelog.NewChangelogFile(dir, marker)

	err := file.Write(want)

	require.NoError(t, err)

	got, err := file.Read()

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestChangelogFileReadError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	marker := "<!-- Fake Marker -->"
	file := changelog.NewChangelogFile(dir, marker)

	got, err := file.Read()

	require.Error(t, err)
	assert.Empty(t, got)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestChangelogFileAddRelease(t *testing.T) {
	tests := []struct {
		name        string
		fixture     string
		marker      string
		wantFixture string
		wantErr     error
		wantErrMsg  string
	}{
		{
			name:    "missing file",
			marker:  "<!-- Fake Marker -->",
			wantErr: os.ErrNotExist,
		},
		{
			name:        "marker only",
			fixture:     "changelog_marker_only.md",
			marker:      "<!-- Fake Marker -->",
			wantFixture: "changelog_marker_only.want.md",
		},
		{
			name:        "existing release",
			fixture:     "changelog.md",
			marker:      "<!-- Fake Marker -->",
			wantFixture: "changelog.want.md",
		},
		{
			name:       "unsupported marker",
			fixture:    "changelog.md",
			marker:     "<!-- Unsupported Marker -->",
			wantErrMsg: `Marker "<!-- Unsupported Marker -->" is missing in CHANGELOG.md`,
		},
		{
			name:       "missing marker",
			fixture:    "changelog_missing_marker.md",
			marker:     "<!-- Fake Marker -->",
			wantErrMsg: `Marker "<!-- Fake Marker -->" is missing in CHANGELOG.md`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()

			if tt.fixture != "" {
				fakeData, err := os.ReadFile(filepath.Join("testdata", tt.fixture))
				require.NoError(t, err)

				err = os.WriteFile(filepath.Join(dir, "CHANGELOG.md"), fakeData, 0644)
				require.NoError(t, err)
			}

			file := changelog.NewChangelogFile(dir, tt.marker)
			release := newMarkdownableFake("Fake Release")

			got, err := file.AddRelease(release)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Empty(t, got)

				return
			}

			if tt.wantErrMsg != "" {
				require.EqualError(t, err, tt.wantErrMsg)
				require.Empty(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, release.Markdown(), got)

			if tt.wantFixture == "" {
				return
			}

			log, err := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
			require.NoError(t, err)

			wantFixtureData, err := os.ReadFile(filepath.Join("testdata", tt.wantFixture))
			require.NoError(t, err)

			assert.Equal(t, string(wantFixtureData), string(log))
		})
	}
}
