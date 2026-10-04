package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const configFilename = ".clg.yml"

func Load(userHomeDir, projectDir string) (Config, error) {
	v := viper.New()

	v.SetDefault("marker", "<!-- CLG -->")
	v.SetDefault("timezone", "UTC")
	v.SetDefault("types", defaultTypes())
	v.SetDefault("markdown.listStyle", "-")

	v.SetEnvPrefix("clg")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Unmarshal only sees known keys, even when AutomaticEnv is enabled.
	for _, key := range []string{
		"marker",
		"timezone",
		"author",
		"markdown.listStyle",
		"markdown.groupsAsList",
	} {
		if err := v.BindEnv(key); err != nil {
			return Config{}, fmt.Errorf("bind environment variable for %q: %w", key, err)
		}
	}

	if err := mergeConfigFile(v, filepath.Join(userHomeDir, configFilename)); err != nil {
		return Config{}, err
	}

	if err := mergeConfigFile(v, filepath.Join(projectDir, configFilename)); err != nil {
		return Config{}, err
	}

	// An explicitly empty author disables attribution; other empty env values stay ignored.
	if author, ok := os.LookupEnv("CLG_AUTHOR"); ok {
		v.Set("author", author)
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
