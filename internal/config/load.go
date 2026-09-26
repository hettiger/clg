package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const configFilename = ".clg.yml"

func Load(userHomeDir, projectDir string) (Config, error) {
	v := viper.New()
	v.SetDefault("marker", "<!-- CLG -->")
	v.SetDefault("types", defaultTypes())
	v.SetDefault("markdown.listStyle", "-")

	if err := mergeConfigFile(v, filepath.Join(userHomeDir, configFilename)); err != nil {
		return Config{}, err
	}

	if err := mergeConfigFile(v, filepath.Join(projectDir, configFilename)); err != nil {
		return Config{}, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func mergeConfigFile(v *viper.Viper, path string) error {
	v.SetConfigFile(path)

	if err := v.MergeInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) || errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	return nil
}
