package changelog

type Release struct {
	Tag    string
	Groups map[Type][]Entry
}

func ReleaseFromUnreleasedEntries(tag string) (Release, error) {
	unreleasedEntries, err := UnreleasedEntries()
	if err != nil {
		return Release{}, err
	}

	supportedTypes := SupportedTypes()
	release := Release{
		Tag:    tag,
		Groups: make(map[Type][]Entry, len(supportedTypes)),
	}
	for _, supportedType := range supportedTypes {
		release.Groups[supportedType] = make([]Entry, 0)
	}

	for _, entry := range unreleasedEntries {
		t, err := TypeFromKeyword(entry.Type)
		if err != nil {
			return Release{}, err
		}
		release.Groups[t] = append(release.Groups[t], entry)
	}

	return release, nil
}
