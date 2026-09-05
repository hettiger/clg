package changelog_test

import (
	"testing"
	"time"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/stretchr/testify/assert"
)

func TestEntry_FilenameAt(t *testing.T) {
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
