package changelog

import (
	"fmt"
	"time"

	"github.com/hettiger/clg/internal/validation"
	"go.yaml.in/yaml/v3"
)

type ChangelogEntryFile struct {
	Entry ChangelogEntry
	Path  string
}

type ChangelogEntry struct {
	Title  string `yaml:"title"`
	Type   string `yaml:"type"`
	Author string `yaml:"author"`
	Group  string `yaml:"group"`
}

func NewChangelogEntry(YAMLData []byte) (ChangelogEntry, error) {
	var entry ChangelogEntry
	if err := yaml.Unmarshal(YAMLData, &entry); err != nil {
		return entry, err
	}
	if err := entry.Validate(); err != nil {
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

func (e ChangelogEntry) Validate() error {
	if err := validation.ValidateMin("title", e.Title, 1); err != nil {
		return err
	}

	if err := validation.ValidateIn("type", e.Type, SupportedTypeKeywords()...); err != nil {
		return err
	}

	return nil
}

func (e ChangelogEntry) Filename(t time.Time) string {
	return fmt.Sprintf(
		"%s-%s.yml",
		t.Format("2006-01-02-150405"),
		e.Type,
	)
}
