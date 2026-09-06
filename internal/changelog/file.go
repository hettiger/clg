package changelog

import (
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Path   string
	Marker string
}

func FileFromCwd(marker string) (File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return File{}, err
	}

	return FileFromDir(cwd, marker), nil
}

func FileFromDir(dir string, marker string) File {
	return File{
		Path:   filepath.Join(dir, "CHANGELOG.md"),
		Marker: marker,
	}
}

func (f File) AddRelease(release Release) (string, error) {
	log, err := f.Read()
	if err != nil {
		return "", err
	}

	md := release.Markdown()
	log = strings.Replace(log, f.Marker, f.Marker+"\n\n"+md+"\n", 1)

	return md, f.Write(log)
}

func (f File) Read() (string, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (f File) Write(log string) error {
	return os.WriteFile(f.Path, []byte(log), 0644)
}
