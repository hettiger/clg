package changelog

import (
	"fmt"
	"strings"
	"time"
)

type Release struct {
	tag      string
	time     time.Time
	sections []section
}

func NewRelease(
	tag string,
	unreleasedEntries []ChangelogEntry,
	releaseTime time.Time,
	groups map[string]string,
	types map[string]string,
) (Release, error) {
	sections := buildSections(groups, types)

	if err := validateSections(sections); err != nil {
		return Release{}, err
	}

	release := Release{
		tag:      tag,
		time:     releaseTime,
		sections: sections,
	}

	for _, entry := range unreleasedEntries {
		if err := addEntryToMatchingSection(release.sections, entry); err != nil {
			return Release{}, err
		}
	}

	return release, nil
}

func (r Release) Markdown() string {
	var result strings.Builder

	fmt.Fprintf(&result, "## [%s] - %s", r.tag, r.time.Format("2006-01-02"))

	switch r.sections[0].kind {
	case groupSectionKind:
		renderGroupSectionsMarkdown(r.sections, &result)

	case typeSectionKind:
		renderTypeSectionsMarkdown(r.sections, "###", &result)
	}

	return result.String()
}

func renderGroupSectionsMarkdown(sections []section, result *strings.Builder) {
	for _, group := range sections {
		if entriesCount(group.children) == 0 {
			continue
		}

		fmt.Fprintf(result, "\n\n### %s", group.headline)

		renderTypeSectionsMarkdown(group.children, "####", result)
	}
}

func renderTypeSectionsMarkdown(
	sections []section,
	headingPrefix string,
	result *strings.Builder,
) {
	for _, section := range sections {
		groupLabel := section.headline
		groupedEntries := section.entries

		if len(groupedEntries) == 0 {
			continue
		}

		groupCount := len(groupedEntries)
		groupCountSuffix := "change"
		if groupCount > 1 {
			groupCountSuffix = "changes"
		}

		fmt.Fprintf(
			result,
			"\n\n%s %s (%d %s)\n",
			headingPrefix,
			groupLabel,
			groupCount,
			groupCountSuffix,
		)

		for _, groupEntry := range groupedEntries {
			fmt.Fprintf(result, "\n- %s", groupEntry.Title)
		}
	}
}

func entriesCount(sections []section) int {
	count := 0

	for _, section := range sections {
		count += len(section.entries)
	}

	return count
}
