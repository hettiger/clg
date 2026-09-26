package changelog

import (
	"testing"

	"github.com/hettiger/clg/internal/testdata"
	"github.com/stretchr/testify/require"
)

func TestSectionKindString(t *testing.T) {
	tests := []struct {
		name string
		kind sectionKind
		want string
	}{
		{
			name: "group",
			kind: 0,
			want: "group",
		},
		{
			name: "type",
			kind: 1,
			want: "type",
		},
		{
			name: "unknown",
			kind: 3,
			want: "unknown (3)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.kind.String()

			require.Equal(t, tt.want, got)
		})
	}
}

func TestBuildSections(t *testing.T) {
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
	section.children = append(section.children, children...)
	return section
}

func typeSection(keyword, headline string, entries ...ChangelogEntry) section {
	section := newTypeSection(keyword, headline)
	section.entries = append(section.entries, entries...)
	return section
}

func TestAddEntryToMatchingSection(t *testing.T) {
	tests := []struct {
		name       string
		sections   []section         // optional: if left empty it is built via groups and types
		groups     map[string]string // used to build sections if empty
		types      map[string]string // used to build sections if empty
		entry      ChangelogEntry
		want       []section
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:   "valid with group",
			groups: testdata.Groups(),
			types:  testdata.Types(),
			entry: ChangelogEntry{
				Group:  "front",
				Type:   "added",
				Title:  "Fake Title",
				Author: "Fake Author",
				Branch: "fake-branch",
			},
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
						typeSection("added", "New Feature", ChangelogEntry{
							Group:  "front",
							Type:   "added",
							Title:  "Fake Title",
							Author: "Fake Author",
							Branch: "fake-branch",
						}),
						typeSection("changed", "Feature Change"),
						typeSection("fixed", "Bug Fix"),
					},
				),
			},
		},
		{
			name:  "valid without group",
			types: testdata.Types(),
			entry: ChangelogEntry{
				Group:  "",
				Type:   "added",
				Title:  "Fake Title",
				Author: "Fake Author",
				Branch: "fake-branch",
			},
			want: []section{
				typeSection("added", "New Feature", ChangelogEntry{
					Group:  "",
					Type:   "added",
					Title:  "Fake Title",
					Author: "Fake Author",
					Branch: "fake-branch",
				}),
				typeSection("changed", "Feature Change"),
				typeSection("fixed", "Bug Fix"),
			},
		},
		{
			name:       "empty sections",
			wantErr:    true,
			wantErrMsg: "no sections available",
		},
		{
			name:  "unsupported group",
			types: testdata.Types(),
			entry: ChangelogEntry{
				Group:  "front",
				Type:   "added",
				Title:  "Fake Title",
				Author: "Fake Author",
				Branch: "fake-branch",
			},
			wantErr:    true,
			wantErrMsg: `entry has group "front", but groups are not configured`,
		},
		{
			name:   "unknown group",
			groups: testdata.Groups(),
			types:  testdata.Types(),
			entry: ChangelogEntry{
				Group:  "weekend",
				Type:   "added",
				Title:  "Fake Title",
				Author: "Fake Author",
				Branch: "fake-branch",
			},
			wantErr:    true,
			wantErrMsg: `unknown entry group: "weekend"`,
		},
		{
			name:  "unknown type",
			types: testdata.Types(),
			entry: ChangelogEntry{
				Group:  "",
				Type:   "special",
				Title:  "Fake Title",
				Author: "Fake Author",
				Branch: "fake-branch",
			},
			wantErr:    true,
			wantErrMsg: `unknown entry type: "special"`,
		},
		{
			name: "mixed section kinds",
			sections: []section{
				section{
					kind: groupSectionKind,
				},
				section{
					kind: typeSectionKind,
				},
			},
			wantErr:    true,
			wantErrMsg: "sections must have the same kind",
		},
		{
			name: "nested mixed section kinds",
			sections: []section{
				section{
					kind:    groupSectionKind,
					keyword: "front",
					children: []section{
						section{
							kind:    groupSectionKind,
							keyword: "back",
						},
						section{
							kind:    typeSectionKind,
							keyword: "changed",
						},
					},
				},
			},
			groups: testdata.Groups(),
			types:  testdata.Types(),
			entry: ChangelogEntry{
				Group:  "front",
				Type:   "changed",
				Title:  "Fake Title",
				Author: "Fake Author",
				Branch: "fake-branch",
			},
			wantErr:    true,
			wantErrMsg: "sections must have the same kind",
		},
		{
			name: "unsupported section kind",
			sections: []section{
				section{
					kind: 4,
				},
			},
			wantErr:    true,
			wantErrMsg: "unsupported section kind: unknown (4)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sections := tt.sections
			if len(sections) == 0 {
				sections = buildSections(tt.groups, tt.types)
			}

			gotErr := addEntryToMatchingSection(sections, tt.entry)

			if tt.wantErr {
				require.Error(t, gotErr)
				require.EqualError(t, gotErr, tt.wantErrMsg)
				return
			}

			require.NoError(t, gotErr)
			require.Equal(t, tt.want, sections)
		})
	}
}
