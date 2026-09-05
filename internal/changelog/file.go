package changelog

import (
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Path string
}

func FileFromCwd() (File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return File{}, err
	}

	return FileFromDir(cwd), nil
}

func FileFromDir(dir string) File {
	return File{
		Path: filepath.Join(dir, "CHANGELOG.md"),
	}
}

func (f File) AddRelease(release Release) error {
	log, err := f.Read()
	if err != nil {
		return err
	}

	marker := "<!-- CLG -->\n\n"
	log = strings.Replace(log, marker, marker+release.Markdown()+"\n\n", 1)

	return f.Write(log)
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
