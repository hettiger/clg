package changelog_test

import (
	"testing"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/stretchr/testify/assert"
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
			got := changelog.NewChangelogFile(tt.dir, tt.marker)
			assert.Equal(t, tt.want, got)
		})
	}
}
