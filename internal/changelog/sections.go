package changelog

import (
	"errors"
	"fmt"

	"github.com/hettiger/clg/internal/support"
)

type Section struct {
	Kind     SectionKind
	Keyword  string
	Headline string
	Children []Section
	Entries  []ChangelogEntry
}

type SectionKind uint8

const (
	GroupSectionKind SectionKind = iota
	TypeSectionKind
)

func (k SectionKind) String() string {
	switch k {
	case GroupSectionKind:
		return "group"

	case TypeSectionKind:
		return "type"

	default:
		return fmt.Sprintf("unknown (%d)", k)
	}
}

func Sections(groups, types map[string]string) []Section {
	sections := groupSections(groups)

	if len(sections) == 0 {
		return typeSections(types)
	}

	for i := range sections {
		sections[i].Children = typeSections(types)
	}

	return sections
}

func groupSections(groups map[string]string) []Section {
	groupKeys := support.SortedMapKeys(groups)
	sections := make([]Section, len(groupKeys))

	for i, groupKey := range groupKeys {
		groupLabel := groups[groupKey]
		sections[i] = Section{
			Kind:     GroupSectionKind,
			Keyword:  groupKey,
			Headline: groupLabel,
			Children: make([]Section, 0),
		}
	}

	return sections
}

func typeSections(types map[string]string) []Section {
	typeKeys := support.SortedMapKeys(types)
	sections := make([]Section, len(typeKeys))

	for i, typeKey := range typeKeys {
		typeLabel := types[typeKey]
		sections[i] = Section{
			Kind:     TypeSectionKind,
			Keyword:  typeKey,
			Headline: typeLabel,
			Entries:  make([]ChangelogEntry, 0),
		}
	}

	return sections
}

func AddEntry(sections []Section, entry ChangelogEntry) error {
	if err := validateSections(sections); err != nil {
		return err
	}

	switch sections[0].Kind {
	case TypeSectionKind:
		if entry.Group != "" {
			return fmt.Errorf("entry has group %q, but groups are not configured", entry.Group)
		}

		return addEntryToTypeSection(sections, entry)

	case GroupSectionKind:
		for i := range sections {
			section := &sections[i]

			if section.Keyword != entry.Group {
				continue
			}

			return addEntryToTypeSection(section.Children, entry)
		}

		return fmt.Errorf("unknown entry group: %q", entry.Group)

	default:
		return fmt.Errorf("unsupported section kind: %s", sections[0].Kind)
	}
}

func addEntryToTypeSection(sections []Section, entry ChangelogEntry) error {
	if err := validateSections(sections); err != nil {
		return err
	}

	if sections[0].Kind != TypeSectionKind {
		return fmt.Errorf("unexpected section kind: %s", sections[0].Kind)
	}

	for i := range sections {
		section := &sections[i]

		if section.Keyword != entry.Type {
			continue
		}

		section.Entries = append(section.Entries, entry)

		return nil
	}

	return fmt.Errorf("unknown entry type: %q", entry.Type)
}

func validateSections(sections []Section) error {
	if len(sections) == 0 {
		return errors.New("no sections available")
	}

	if !sectionsAreUniform(sections) {
		return errors.New("sections of mixed kinds are not supported")
	}

	switch sections[0].Kind {
	case GroupSectionKind, TypeSectionKind:
		return nil

	default:
		return fmt.Errorf("unsupported section kind: %s", sections[0].Kind)
	}
}

func sectionsAreUniform(sections []Section) bool {
	if len(sections) == 0 {
		return true
	}

	kind := sections[0].Kind

	for _, section := range sections[1:] {
		if section.Kind != kind {
			return false
		}
	}

	return true
}
