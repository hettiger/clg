package changelog

import (
	"fmt"
	"strings"
	"time"
)

type Release struct {
	Tag    string
	Groups map[Type][]Entry
}

func ReleaseFromUnreleasedEntries(tag string) (Release, error) {
	unreleasedEntries, err := UnreleasedEntries()
	if err != nil {
		return Release{}, err
	}

	supportedTypes := SupportedTypes()
	release := Release{
		Tag:    tag,
		Groups: make(map[Type][]Entry, len(supportedTypes)),
	}
	for _, supportedType := range supportedTypes {
		release.Groups[supportedType] = make([]Entry, 0)
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

	fmt.Fprintf(&result, "## [%s] - %s", r.Tag, time.Now().UTC().Format("2006-01-02"))

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
