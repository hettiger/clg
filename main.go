package main

import (
	"log"
	"os"
	"time"

	"github.com/hettiger/clg/cmd"
	"github.com/hettiger/clg/internal/changelog"
	"github.com/hettiger/clg/internal/config"
	"github.com/spf13/viper"
)

func main() {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var config config.Config
	viper.SetConfigName("clg")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(rootDir)
	viper.SetDefault("Marker", "<!-- CLG -->")
	viper.SetDefault("Types", map[string]string{
		"added":       "New Feature",
		"fixed":       "Bug Fix",
		"hotfix":      "Hotfix",
		"changed":     "Feature Change",
		"deprecated":  "New Deprecation",
		"removed":     "Feature Removal",
		"security":    "Security Fix",
		"performance": "Performance Improvement",
		"other":       "Other",
		"ignore":      "No Changelog",
	})
	viper.SetDefault("Markdown.ListStyle", "-")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err)
		}
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	now := func() time.Time {
		return time.Now().UTC()
	}

	entryStore := changelog.NewEntryStore(rootDir, now)

	app := cmd.NewApp(config, now, rootDir, entryStore)

	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
