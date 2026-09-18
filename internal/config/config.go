package config

type Config struct {
	Marker   string         `mapstructure:"marker"`
	Groups   []string       `mapstructure:"groups"`
	Types    Types          `mapstructure:"types"`
	Markdown MarkdownConfig `mapstructure:"markdown"`
}

type Types map[string]string

type MarkdownConfig struct {
	ListStyle    string `mapstructure:"listStyle"`
	GroupsAsList bool   `mapstructure:"groupsAsList"`
}
