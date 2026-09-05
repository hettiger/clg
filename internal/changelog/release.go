package changelog

import (
	"fmt"
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

func (r Release) Markdown() (string, error) {
	result := ""

	result += fmt.Sprintf("## [%s] - %s", r.Tag, time.Now().UTC().Format("2006-01-02"))
	result += fmt.Sprintln()

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

		result += fmt.Sprintln()
		result += fmt.Sprintf("### %s (%v %s)", groupType.Label, groupCount, groupCountSuffix)
		result += fmt.Sprintln()
		result += fmt.Sprintln()

		for _, groupEntry := range groupedEntries {
			result += fmt.Sprintln("- " + groupEntry.Title)
		}
	}

	return result, nil
}
