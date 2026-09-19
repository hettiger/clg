package changelog

import (
	"testing"

	"github.com/hettiger/clg/internal/testdata"
	"github.com/stretchr/testify/require"
)

func TestSections(t *testing.T) {
	tests := []struct {
		name   string
		groups map[string]string
		types  map[string]string
		want   []section
	}{
		{
			name:  "without groups",
			types: testdata.Types(),
			want: []section{
				typeSection("added", "New Feature"),
				typeSection("changed", "Feature Change"),
				typeSection("fixed", "Bug Fix"),
			},
		},
		{
			name:   "with groups",
			groups: testdata.Groups(),
			types:  testdata.Types(),
			want: []section{
				groupSection(
					"back",
					"Backend",
					[]section{
						typeSection("added", "New Feature"),
						typeSection("changed", "Feature Change"),
						typeSection("fixed", "Bug Fix"),
					},
				),
				groupSection(
					"front",
					"Frontend",
					[]section{
						typeSection("added", "New Feature"),
						typeSection("changed", "Feature Change"),
						typeSection("fixed", "Bug Fix"),
					},
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := buildSections(tt.groups, tt.types)

			require.Equal(t, tt.want, got)
		})
	}
}

func groupSection(keyword, headline string, children []section) section {
	section := newGroupSection(keyword, headline)
	section.children = children
	return section
}

var typeSection = newTypeSection
