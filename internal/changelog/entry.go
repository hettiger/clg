package changelog

import (
	"fmt"
	"strings"

	"github.com/hettiger/clg/internal/validation"
	"go.yaml.in/yaml/v3"
)

type ChangelogEntryFile struct {
	Entry ChangelogEntry
	Path  string
}

type ChangelogEntry struct {
	Group  string `yaml:"group"`
	Type   string `yaml:"type"`
	Title  string `yaml:"title"`
	Author string `yaml:"author"`
	Branch string `yaml:"branch"`
}

func NewChangelogEntry(YAMLData []byte, groupKeys, typeKeys []string) (ChangelogEntry, error) {
	var entry ChangelogEntry
	if err := yaml.Unmarshal(YAMLData, &entry); err != nil {
		return entry, err
	}
	if err := entry.Validate(groupKeys, typeKeys); err != nil {
		return ChangelogEntry{}, err
	}
	return entry, nil
}

func (e ChangelogEntry) YAMLData() ([]byte, error) {
	return yaml.Marshal(e)
}

func (e ChangelogEntry) YAML() (string, error) {
	data, err := e.YAMLData()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (e ChangelogEntry) Validate(groupKeys, typeKeys []string) error {
	if err := validation.ValidateMin("title", e.Title, 1); err != nil {
		return err
	}

	if len(groupKeys) == 0 {
		if e.Group != "" {
			return fmt.Errorf("Unsupported group (%s)", e.Group)
		}
	} else if err := validation.ValidateIn(
		"group",
		e.Group,
		groupKeys...,
	); err != nil {
		return err
	}

	if err := validation.ValidateIn(
		"type",
		e.Type,
		typeKeys...,
	); err != nil {
		return err
	}

	return nil
}

func (e ChangelogEntry) Filename(uuidV7 func() (string, error)) (string, error) {
	id, err := uuidV7()
	if err != nil {
		return "", err
	}

	var filename strings.Builder

	if e.Group != "" {
		fmt.Fprintf(&filename, "%s-", e.Group)
	}

	fmt.Fprintf(&filename, "%s-%s.yml", e.Type, id)

	return filename.String(), nil
}
