package main

import (
	"encoding/json"
	"fmt"
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
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	jsonData, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))

	now := func() time.Time {
		return time.Now().UTC()
	}

	entryStore := changelog.NewEntryStore(rootDir, now)

	app := cmd.NewApp(config, now, rootDir, entryStore)

	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
