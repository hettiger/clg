package config

type Config struct {
	Marker   string            `mapstructure:"marker"`
	Groups   map[string]string `mapstructure:"groups"`
	Types    map[string]string `mapstructure:"types"`
	Author   *string           `mapstructure:"author"`
	Markdown MarkdownConfig    `mapstructure:"markdown"`
	Issue    IssueConfig       `mapstructure:"issue"`
}

type MarkdownConfig struct {
	ListStyle    string `mapstructure:"listStyle"`
	GroupsAsList bool   `mapstructure:"groupsAsList"`
}

type IssueConfig struct {
	DisplayPrefix string `mapstructure:"displayPrefix"`
	Pattern       string `mapstructure:"pattern"`
}
