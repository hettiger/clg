package config

type Config struct {
	Marker   string            `mapstructure:"marker"`
	Timezone string            `mapstructure:"timezone"`
	Groups   map[string]string `mapstructure:"groups"`
	Types    map[string]string `mapstructure:"types"`
	Author   *string           `mapstructure:"author"`
	Markdown MarkdownConfig    `mapstructure:"markdown"`
}

type MarkdownConfig struct {
	ListStyle    string `mapstructure:"listStyle"`
	GroupsAsList bool   `mapstructure:"groupsAsList"`
}
