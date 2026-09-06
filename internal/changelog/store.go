package changelog

import (
	"os"
	"path/filepath"
	"time"
)

type EntryStore struct {
	rootDir string
	now     func() time.Time
}

func NewEntryStore(root string, now func() time.Time) EntryStore {
	return EntryStore{
		rootDir: root,
		now:     now,
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

		entry, err := NewChangelogEntry(data)
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
	data, err := entry.YAMLData()
	if err != nil {
		return "", err
	}

	dir, err := s.unreleasedDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, entry.Filename(s.now()))

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
