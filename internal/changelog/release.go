package changelog

import (
	"fmt"
	"maps"
	"strings"
	"time"

	"github.com/hettiger/clg/internal/support"
)

type Release struct {
	tag      string
	time     time.Time
	groups   map[string][]ChangelogEntry
	types    map[string]string
	typeKeys []string
}

func NewRelease(tag string, unreleasedEntries []ChangelogEntry, releaseTime time.Time, types map[string]string) (Release, error) {
	types = maps.Clone(types)

	release := Release{
		tag:      tag,
		time:     releaseTime,
		groups:   make(map[string][]ChangelogEntry, len(types)),
		types:    types,
		typeKeys: support.SortedMapKeys(types),
	}
	for _, typeKey := range release.typeKeys {
		release.groups[typeKey] = make([]ChangelogEntry, 0)
	}

	for _, entry := range unreleasedEntries {
		if _, ok := types[entry.Type]; !ok {
			return Release{}, fmt.Errorf("Unknown type keyword provided (%s)", entry.Type)
		}

		release.groups[entry.Type] = append(release.groups[entry.Type], entry)
	}

	return release, nil
}

func (r Release) Markdown() string {
	var result strings.Builder

	fmt.Fprintf(&result, "## [%s] - %s", r.tag, r.time.Format("2006-01-02"))

	for _, typeKey := range r.typeKeys {
		if typeKey == "ignore" {
			continue
		}

		groupLabel := r.types[typeKey]
		groupedEntries := r.groups[typeKey]

		if len(groupedEntries) == 0 {
			continue
		}

		groupCount := len(groupedEntries)
		groupCountSuffix := "change"
		if groupCount > 1 {
			groupCountSuffix = "changes"
		}

		fmt.Fprintf(&result, "\n\n### %s (%d %s)\n", groupLabel, groupCount, groupCountSuffix)

		for _, groupEntry := range groupedEntries {
			fmt.Fprintf(&result, "\n- %s", groupEntry.Title)
		}
	}

	return result.String()
}
