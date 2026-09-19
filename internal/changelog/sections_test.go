package changelog_test

import (
	"testing"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestSections(t *testing.T) {
	tests := []struct {
		name   string
		groups map[string]string
		types  map[string]string
		want   []changelog.Section
	}{
		{
			name:  "without groups",
			types: testTypes(),
			want: []changelog.Section{
				typeSection("added", "New Feature"),
				typeSection("changed", "Feature Change"),
				typeSection("fixed", "Bug Fix"),
			},
		},
		{
			name:   "with groups",
			groups: testGroups(),
			types:  testTypes(),
			want: []changelog.Section{
				groupSection(
					"back",
					"Backend",
					[]changelog.Section{
						typeSection("added", "New Feature"),
						typeSection("changed", "Feature Change"),
						typeSection("fixed", "Bug Fix"),
					},
				),
				groupSection(
					"front",
					"Frontend",
					[]changelog.Section{
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

			got := changelog.Sections(tt.groups, tt.types)

			require.Equal(t, tt.want, got)
		})
	}
}

func groupSection(keyword, headline string, children []changelog.Section) changelog.Section {
	return changelog.Section{
		Kind:     changelog.GroupSectionKind,
		Keyword:  keyword,
		Headline: headline,
		Children: children,
	}
}

func typeSection(keyword, headline string) changelog.Section {
	return changelog.Section{
		Kind:     changelog.TypeSectionKind,
		Keyword:  keyword,
		Headline: headline,
		Entries:  make([]changelog.ChangelogEntry, 0),
	}
}
