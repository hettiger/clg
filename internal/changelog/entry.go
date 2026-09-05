package changelog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.yaml.in/yaml/v3"
)

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
	dir, err := unreleasedDir()
	if err != nil {
		return nil, err
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(dirEntries))

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

		entries = append(entries, entry)
	}

	return entries, nil
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
