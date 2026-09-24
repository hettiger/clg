package changelog

import (
	"os"
	"path/filepath"
	"time"

	"github.com/hettiger/clg/internal/support"
)

type EntryStore struct {
	rootDir   string
	now       func() time.Time
	uuidV7    func() (string, error)
	groupKeys []string
	typeKeys  []string
}

func NewEntryStore(
	root string,
	now func() time.Time,
	uuidV7 func() (string, error),
	groups,
	types map[string]string,
) EntryStore {
	return EntryStore{
		rootDir:   root,
		now:       now,
		uuidV7:    uuidV7,
		groupKeys: support.SortedMapKeys(groups),
		typeKeys:  support.SortedMapKeys(types),
	}
}

func (s EntryStore) UnreleasedEntries() ([]ChangelogEntry, error) {
	entryFiles, err := s.UnreleasedEntryFiles()
	if err != nil {
		return nil, err
	}

	entries := make([]ChangelogEntry, len(entryFiles))
	for i, ef := range entryFiles {
		entries[i] = ef.Entry
	}

	return entries, nil
}

func (s EntryStore) UnreleasedEntryFiles() ([]ChangelogEntryFile, error) {
	dir, err := s.unreleasedDir()
	if err != nil {
		return nil, err
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	entryFiles := make([]ChangelogEntryFile, 0, len(dirEntries))

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() || filepath.Ext(dirEntry.Name()) != ".yml" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, dirEntry.Name()))
		if err != nil {
			return nil, err
		}

		entry, err := NewChangelogEntry(data, s.groupKeys, s.typeKeys)
		if err != nil {
			return nil, err
		}

		entryFiles = append(entryFiles, ChangelogEntryFile{
			Entry: entry,
			Path:  filepath.Join(dir, dirEntry.Name()),
		})
	}

	return entryFiles, nil
}

func (s EntryStore) Write(entry ChangelogEntry) (string, error) {
	if err := entry.Validate(s.groupKeys, s.typeKeys); err != nil {
		return "", err
	}

	data, err := entry.YAMLData()
	if err != nil {
		return "", err
	}

	dir, err := s.unreleasedDir()
	if err != nil {
		return "", err
	}

	filename, err := entry.Filename(s.uuidV7)
	if err != nil {
		return "", nil
	}

	path := filepath.Join(dir, filename)

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}

	return path, nil
}

func (s EntryStore) unreleasedDir() (string, error) {
	pathDir := filepath.Join(s.rootDir, "changelogs", "unreleased")

	if err := os.MkdirAll(pathDir, 0755); err != nil {
		return "", err
	}

	return pathDir, nil
}
