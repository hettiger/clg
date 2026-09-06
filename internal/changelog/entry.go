package changelog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hettiger/clg/internal/validation"
	"go.yaml.in/yaml/v3"
)

type EntryFile struct {
	Entry Entry
	Path  string
}

type Entry struct {
	Title  string `yaml:"title"`
	Type   string `yaml:"type"`
	Author string `yaml:"author"`
	Group  string `yaml:"group"`
}

func EntryFromYAML(data []byte) (Entry, error) {
	var entry Entry
	if err := yaml.Unmarshal(data, &entry); err != nil {
		return entry, err
	}
	if err := entry.Validate(); err != nil {
		return Entry{}, err
	}
	return entry, nil
}

func (e Entry) YAMLData() ([]byte, error) {
	return yaml.Marshal(e)
}

func (e Entry) YAML() (string, error) {
	data, err := e.YAMLData()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (e Entry) Validate() error {
	if err := validation.ValidateMin("title", e.Title, 1); err != nil {
		return err
	}

	if err := validation.ValidateIn("type", e.Type, SupportedTypeKeywords()...); err != nil {
		return err
	}

	return nil
}

func (e Entry) Filename() string {
	return e.FilenameAt(time.Now())
}

func (e Entry) FilenameAt(t time.Time) string {
	return fmt.Sprintf(
		"%s-%s.yml",
		t.UTC().Format("2006-01-02-150405"),
		e.Type,
	)
}

func (e Entry) WriteToFile() error {
	data, err := e.YAMLData()
	if err != nil {
		return err
	}

	dir, err := unreleasedDir()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, e.Filename()), data, 0644)
}

func UnreleasedEntries() ([]Entry, error) {
	entryFiles, err := UnreleasedEntryFiles()
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, len(entryFiles))
	for i, ef := range entryFiles {
		entries[i] = ef.Entry
	}

	return entries, nil
}

func UnreleasedEntryFiles() ([]EntryFile, error) {
	dir, err := unreleasedDir()
	if err != nil {
		return nil, err
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	entryFiles := make([]EntryFile, 0, len(dirEntries))

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() || filepath.Ext(dirEntry.Name()) != ".yml" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, dirEntry.Name()))
		if err != nil {
			return nil, err
		}

		entry, err := EntryFromYAML(data)
		if err != nil {
			return nil, err
		}

		entryFiles = append(entryFiles, EntryFile{
			Entry: entry,
			Path:  filepath.Join(dir, dirEntry.Name()),
		})
	}

	return entryFiles, nil
}

func unreleasedDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	pathDir := filepath.Join(cwd, "changelogs", "unreleased")

	if err := os.MkdirAll(pathDir, 0755); err != nil {
		return "", err
	}

	return pathDir, nil
}
