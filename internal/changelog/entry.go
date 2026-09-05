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

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	pathDir := filepath.Join(cwd, "changelogs", "unreleased")
	pathFile := filepath.Join(pathDir, e.Filename())

	if err := os.MkdirAll(pathDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(pathFile, data, 0644)
}
