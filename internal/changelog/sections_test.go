package changelog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hettiger/clg/internal/support"
	"github.com/hettiger/clg/internal/testdata"
	"github.com/stretchr/testify/require"
)

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
		name         string
		sections     []section         // optional: if left empty it is built via groups and types
		groups       map[string]string // used to build sections if empty and to validate ChangelogEntry if entryFixture is provided
		types        map[string]string // used to build sections if empty and to validate ChangelogEntry if entryFixture is provided
		entry        ChangelogEntry    // optional: use entry or entryFixture
		entryFixture string            // optional: use entry or entryFixture
		want         []section
		wantErr      bool
		wantErrMsg   string
	}{
		{
			name:         "valid with group",
			groups:       testdata.Groups(),
			types:        testdata.Types(),
			entryFixture: "entry_valid_group.yml",
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
							Title:  "Fake Title",
							Type:   "added",
							Author: "Fake Author",
							Group:  "front",
						}),
						typeSection("changed", "Feature Change"),
						typeSection("fixed", "Bug Fix"),
					},
				),
			},
		},
		{
			name:         "valid without group",
			types:        testdata.Types(),
			entryFixture: "entry_valid.yml",
			want: []section{
				typeSection("added", "New Feature", ChangelogEntry{
					Title:  "Fake Title",
					Type:   "added",
					Author: "Fake Author",
					Group:  "",
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
				Title:  "Fake Title",
				Type:   "added",
				Author: "Fake Author",
				Group:  "front",
			},
			wantErr:    true,
			wantErrMsg: `entry has group "front", but groups are not configured`,
		},
		{
			name:   "unknown group",
			groups: testdata.Groups(),
			types:  testdata.Types(),
			entry: ChangelogEntry{
				Title:  "Fake Title",
				Type:   "added",
				Author: "Fake Author",
				Group:  "weekend",
			},
			wantErr:    true,
			wantErrMsg: `unknown entry group: "weekend"`,
		},
		{
			name:  "unknown type",
			types: testdata.Types(),
			entry: ChangelogEntry{
				Title:  "Fake Title",
				Type:   "special",
				Author: "Fake Author",
				Group:  "",
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
					kind: groupSectionKind,
					children: []section{
						section{
							kind: groupSectionKind,
						},
						section{
							kind: typeSectionKind,
						},
					},
				},
			},
			wantErr:    true,
			wantErrMsg: "sections must have the same kind",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			entry := tt.entry

			if tt.entryFixture != "" {
				entryFixtureData, err := os.ReadFile(filepath.Join("testdata", tt.entryFixture))
				require.NoError(t, err)
				groupKeys := support.SortedMapKeys(tt.groups)
				typeKeys := support.SortedMapKeys(tt.types)
				entry, err = NewChangelogEntry(entryFixtureData, groupKeys, typeKeys)
				require.NoError(t, err)
			}

			sections := tt.sections

			if len(sections) == 0 {
				sections = buildSections(tt.groups, tt.types)
			}

			gotErr := addEntryToMatchingSection(sections, entry)

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
