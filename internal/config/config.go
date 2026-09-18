package config

type Config struct {
	Marker   string
	Groups   []string
	Types    map[string]string
	Markdown MarkdownConfig
}

type MarkdownConfig struct {
	ListStyle    string
	GroupsAsList bool
}
