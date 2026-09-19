package config

import (
	"github.com/spf13/viper"
)

func Load(rootDir string) (Config, error) {
	var config Config

	v := viper.New()

	v.SetConfigName("clg")
	v.SetConfigType("yaml")
	v.AddConfigPath(rootDir)

	v.SetDefault("marker", "<!-- CLG -->")
	v.SetDefault("types", defaultTypes())
	v.SetDefault("markdown.listStyle", "-")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	if err := v.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
