package changelog

import (
	"fmt"
	"strings"
	"time"
)

type Release struct {
	Tag    string
	Groups map[Type][]ChangelogEntry
	time   time.Time
}

func NewRelease(tag string, unreleasedEntries []ChangelogEntry, time time.Time) (Release, error) {
	supportedTypes := SupportedTypes()
	release := Release{
		Tag:    tag,
		Groups: make(map[Type][]ChangelogEntry, len(supportedTypes)),
		time:   time,
	}
	for _, supportedType := range supportedTypes {
		release.Groups[supportedType] = make([]ChangelogEntry, 0)
	}

	for _, entry := range unreleasedEntries {
		t, err := TypeFromKeyword(entry.Type)
		if err != nil {
			return Release{}, err
		}
		release.Groups[t] = append(release.Groups[t], entry)
	}

	return release, nil
}

func (r Release) Markdown() string {
	var result strings.Builder

	fmt.Fprintf(&result, "## [%s] - %s", r.Tag, r.time.Format("2006-01-02"))

	for _, groupType := range SupportedTypes() {
		if groupType.Keyword == "ignore" {
			continue
		}

		groupedEntries := r.Groups[groupType]
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
