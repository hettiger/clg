package changelog

import (
	"errors"
	"fmt"

	"github.com/hettiger/clg/internal/support"
)

type section struct {
	kind     sectionKind
	keyword  string
	headline string
	children []section
	entries  []ChangelogEntry
}

type sectionKind uint8

const (
	groupSectionKind sectionKind = iota
	typeSectionKind
)

func (k sectionKind) String() string {
	switch k {
	case groupSectionKind:
		return "group"

	case typeSectionKind:
		return "type"

	default:
		return fmt.Sprintf("unknown (%d)", k)
	}
}

func newGroupSection(keyword, headline string) section {
	return section{
		kind:     groupSectionKind,
		keyword:  keyword,
		headline: headline,
		children: make([]section, 0),
	}
}

func newTypeSection(keyword, headline string) section {
	return section{
		kind:     typeSectionKind,
		keyword:  keyword,
		headline: headline,
		entries:  make([]ChangelogEntry, 0),
	}
}

func buildSections(groups, types map[string]string) []section {
	sections := buildGroupSections(groups)

	if len(sections) == 0 {
		return buildTypeSections(types)
	}

	for i := range sections {
		sections[i].children = buildTypeSections(types)
	}

	return sections
}

func buildGroupSections(groups map[string]string) []section {
	groupKeys := support.SortedMapKeys(groups)
	sections := make([]section, len(groupKeys))

	for i, groupKey := range groupKeys {
		groupLabel := groups[groupKey]
		sections[i] = newGroupSection(groupKey, groupLabel)
	}

	return sections
}

func buildTypeSections(types map[string]string) []section {
	typeKeys := support.SortedMapKeys(types)
	sections := make([]section, len(typeKeys))

	for i, typeKey := range typeKeys {
		typeLabel := types[typeKey]
		sections[i] = newTypeSection(typeKey, typeLabel)
	}

	return sections
}

func addEntryToMatchingSection(sections []section, entry ChangelogEntry) error {
	if err := validateSections(sections); err != nil {
		return err
	}

	switch sections[0].kind {
	case typeSectionKind:
		if entry.Group != "" {
			return fmt.Errorf("entry has group %q, but groups are not configured", entry.Group)
		}

		return addEntryToMatchingTypeSection(sections, entry)

	case groupSectionKind:
		for i := range sections {
			section := &sections[i]

			if section.keyword != entry.Group {
				continue
			}

			return addEntryToMatchingTypeSection(section.children, entry)
		}

		return fmt.Errorf("unknown entry group: %q", entry.Group)

	default:
		return fmt.Errorf("unsupported section kind: %s", sections[0].kind)
	}
}

func addEntryToMatchingTypeSection(sections []section, entry ChangelogEntry) error {
	if err := validateSections(sections); err != nil {
		return err
	}

	if sections[0].kind != typeSectionKind {
		return fmt.Errorf("unexpected section kind: %s", sections[0].kind)
	}

	for i := range sections {
		section := &sections[i]

		if section.keyword != entry.Type {
			continue
		}

		section.entries = append(section.entries, entry)

		return nil
	}

	return fmt.Errorf("unknown entry type: %q", entry.Type)
}

func validateSections(sections []section) error {
	if len(sections) == 0 {
		return errors.New("no sections available")
	}

	if err := validateUniformSections(sections); err != nil {
		return err
	}

	switch sections[0].kind {
	case groupSectionKind, typeSectionKind:
		return nil

	default:
		return fmt.Errorf("unsupported section kind: %s", sections[0].kind)
	}
}

func validateUniformSections(sections []section) error {
	if len(sections) == 0 {
		return nil
	}

	kind := sections[0].kind

	for _, section := range sections[1:] {
		if section.kind != kind {
			return errors.New("sections must have the same kind")
		}
	}

	return nil
}
