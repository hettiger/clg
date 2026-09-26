package config

type Config struct {
	Marker   string            `mapstructure:"marker"`
	Groups   map[string]string `mapstructure:"groups"`
	Types    map[string]string `mapstructure:"types"`
	Author   AuthorConfig      `mapstructure:"author"`
	Markdown MarkdownConfig    `mapstructure:"markdown"`
}

type AuthorConfig struct {
	Name string `mapstructure:"name"`
	URL  string `mapstructure:"url"`
}

type MarkdownConfig struct {
	ListStyle    string `mapstructure:"listStyle"`
	GroupsAsList bool   `mapstructure:"groupsAsList"`
}
