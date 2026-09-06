package changelog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ChangelogFile struct {
	Path   string
	Marker string
}

func NewChangelogFile(dir string, marker string) ChangelogFile {
	return ChangelogFile{
		Path:   filepath.Join(dir, "CHANGELOG.md"),
		Marker: marker,
	}
}

func (f ChangelogFile) AddRelease(release Release) (string, error) {
	log, err := f.Read()
	if err != nil {
		return "", err
	}

	if !strings.Contains(log, f.Marker) {
		return "", fmt.Errorf(`Marker "%s" is missing in CHANGELOG.md`, f.Marker)
	}

	md := release.Markdown()
	log = strings.Replace(log, f.Marker, f.Marker+"\n\n"+md+"\n", 1)

	return md, f.Write(log)
}

func (f ChangelogFile) Read() (string, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (f ChangelogFile) Write(log string) error {
	return os.WriteFile(f.Path, []byte(log), 0644)
}
