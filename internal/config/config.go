package config

type Config struct {
	Marker   string            `mapstructure:"marker"`
	Groups   map[string]string `mapstructure:"groups"`
	Types    map[string]string `mapstructure:"types"`
	Markdown MarkdownConfig    `mapstructure:"markdown"`
}

type MarkdownConfig struct {
	ListStyle    string `mapstructure:"listStyle"`
	GroupsAsList bool   `mapstructure:"groupsAsList"`
}
