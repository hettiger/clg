package changelog

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Release struct {
	tag                string
	time               time.Time
	sections           []section
	issueDisplayPrefix string
	issuePattern       string
}

func NewRelease(
	tag string,
	unreleasedEntries []ChangelogEntry,
	releaseTime time.Time,
	groups map[string]string,
	types map[string]string,
	issueDisplayPrefix string,
	issuePattern string,
) (Release, error) {
	sections := buildSections(groups, types)

	if err := validateSections(sections); err != nil {
		return Release{}, err
	}

	release := Release{
		tag:                tag,
		time:               releaseTime,
		sections:           sections,
		issueDisplayPrefix: issueDisplayPrefix,
		issuePattern:       issuePattern,
	}

	for _, entry := range unreleasedEntries {
		if err := addEntryToMatchingSection(release.sections, entry); err != nil {
			return Release{}, err
		}
	}

	return release, nil
}

func (r Release) Markdown() (string, error) {
	var result strings.Builder

	if _, err := fmt.Fprintf(&result, "## [%s] - %s", r.tag, r.time.Format("2006-01-02")); err != nil {
		return "", err
	}

	switch r.sections[0].kind {
	case groupSectionKind:
		if err := r.renderGroupSectionsMarkdown(r.sections, &result); err != nil {
			return "", err
		}

	case typeSectionKind:
		if err := r.renderTypeSectionsMarkdown(r.sections, "###", &result); err != nil {
			return "", err
		}
	}

	return result.String(), nil
}

func (r Release) renderGroupSectionsMarkdown(sections []section, result *strings.Builder) error {
	for _, group := range sections {
		if entriesCount(group.children) == 0 {
			continue
		}

		if _, err := fmt.Fprintf(result, "\n\n### %s", group.headline); err != nil {
			return err
		}

		if err := r.renderTypeSectionsMarkdown(group.children, "####", result); err != nil {
			return err
		}
	}

	return nil
}

func (r Release) renderTypeSectionsMarkdown(
	sections []section,
	headingPrefix string,
	result *strings.Builder,
) error {
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

		if _, err := fmt.Fprintf(
			result,
			"\n\n%s %s (%d %s)\n",
			headingPrefix,
			groupLabel,
			groupCount,
			groupCountSuffix,
		); err != nil {
			return err
		}

		for _, e := range groupedEntries {
			issue, err := e.Issue(r.issuePattern)
			switch {
			case err == nil, errors.Is(err, ErrIssuePatternMismatch):
				// continue on success or unsupported branch names
			default:
				return err
			}

			var meta string
			switch {
			case issue != "" && e.Author != "":
				meta = fmt.Sprintf(" (%s%s, %s)", r.issueDisplayPrefix, issue, e.Author)
			case issue != "":
				meta = fmt.Sprintf(" (%s%s)", r.issueDisplayPrefix, issue)
			case e.Author != "":
				meta = fmt.Sprintf(" (%s)", e.Author)
			default:
				meta = ""
			}

			if _, err := fmt.Fprintf(result, "\n- %s%s", e.Title, meta); err != nil {
				return err
			}
		}
	}

	return nil
}

func entriesCount(sections []section) int {
	count := 0

	for _, section := range sections {
		count += len(section.entries)
	}

	return count
}
