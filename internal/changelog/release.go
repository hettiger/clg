package changelog

import (
	"fmt"
	"strings"
	"time"
)

type Release struct {
	tag    string
	groups map[Type][]ChangelogEntry
	time   time.Time
}

func NewRelease(tag string, unreleasedEntries []ChangelogEntry, time time.Time) (Release, error) {
	supportedTypes := SupportedTypes()
	release := Release{
		tag:    tag,
		groups: make(map[Type][]ChangelogEntry, len(supportedTypes)),
		time:   time,
	}
	for _, supportedType := range supportedTypes {
		release.groups[supportedType] = make([]ChangelogEntry, 0)
	}

	for _, entry := range unreleasedEntries {
		t, err := TypeFromKeyword(entry.Type)
		if err != nil {
			return Release{}, err
		}
		release.groups[t] = append(release.groups[t], entry)
	}

	return release, nil
}

func (r Release) Markdown() string {
	var result strings.Builder

	fmt.Fprintf(&result, "## [%s] - %s", r.tag, r.time.Format("2006-01-02"))

	for _, groupType := range SupportedTypes() {
		if groupType.Keyword == "ignore" {
			continue
		}

		groupedEntries := r.groups[groupType]
		if len(groupedEntries) == 0 {
			continue
		}

		groupCount := len(groupedEntries)
		groupCountSuffix := "change"
		if groupCount > 1 {
			groupCountSuffix = "changes"
		}

		fmt.Fprintf(&result, "\n\n### %s (%d %s)\n", groupType.Label, groupCount, groupCountSuffix)

		for _, groupEntry := range groupedEntries {
			fmt.Fprintf(&result, "\n- %s", groupEntry.Title)
		}
	}

	return result.String()
}
