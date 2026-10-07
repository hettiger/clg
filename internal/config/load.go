package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const configFilename = ".clg.yml"

func Load(userHomeDir, projectDir string) (Config, error) {
	v := viper.New()

	v.SetDefault("marker", "<!-- CLG -->")
	v.SetDefault("timezone", "UTC")
	v.SetDefault("types", defaultTypes())
	v.SetDefault("markdown.listStyle", "-")
	v.SetDefault("issue.displayPrefix", "#")
	v.SetDefault("issue.pattern", "^(?:[^\\/]+\\/)?(\\d+)-\\S+$") // extracts `31` in `feature/31-description` or `31-description`

	for _, binding := range []struct {
		key     string
		envVars []string
	}{
		{key: "marker", envVars: []string{"CLG_MARKER"}},
		{key: "timezone", envVars: []string{"CLG_TIMEZONE"}},
		{key: "author", envVars: []string{"CLG_AUTHOR"}},
		{key: "markdown.listStyle", envVars: []string{"CLG_MARKDOWN_LIST_STYLE"}},
		{key: "markdown.groupsAsList", envVars: []string{"CLG_MARKDOWN_GROUPS_AS_LIST"}},
		{key: "issue.displayPrefix", envVars: []string{"CLG_ISSUE_DISPLAY_PREFIX"}},
		{key: "issue.pattern", envVars: []string{"CLG_ISSUE_PATTERN"}},
	} {
		envVars := append([]string{binding.key}, binding.envVars...)
		if err := v.BindEnv(envVars...); err != nil {
			return Config{}, fmt.Errorf("bind environment variable for %q: %w", binding.key, err)
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
