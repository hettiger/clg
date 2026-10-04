package changelog_test

type markdownableFake struct {
	value string
}

func newMarkdownableFake(value string) markdownableFake {
	return markdownableFake{
		value: value,
	}
}

func (f markdownableFake) Markdown() (string, error) {
	return f.value, nil
}
