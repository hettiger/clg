package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/renameio/v2/maybe"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("project config overrides user config", func(t *testing.T) {
		t.Setenv("CLG_TIMEZONE", "")

		homeDir := t.TempDir()
		projectDir := t.TempDir()

		writeConfig(t, homeDir, "timezone: Europe/Paris\n")
		writeConfig(t, projectDir, "timezone: Asia/Kolkata\n")

		cfg, err := Load(homeDir, projectDir)
		require.NoError(t, err)

		require.Equal(t, "Asia/Kolkata", cfg.Timezone)
	})

	t.Run("it binds env vars", func(t *testing.T) {
		t.Setenv("CLG_MARKER", "<!-- from env -->")
		t.Setenv("CLG_TIMEZONE", "Europe/Paris")
		t.Setenv("CLG_AUTHOR", "env author")
		t.Setenv("CLG_MARKDOWN_LIST_STYLE", "*")
		t.Setenv("CLG_MARKDOWN_GROUPS_AS_LIST", "true")
		t.Setenv("CLG_ISSUE_DISPLAY_PREFIX", "issue-")
		t.Setenv("CLG_ISSUE_PATTERN", "^task/(\\d+)-\\S+$")

		cfg, err := Load(t.TempDir(), t.TempDir())
		require.NoError(t, err)

		require.Equal(t, "<!-- from env -->", cfg.Marker)
		require.Equal(t, "Europe/Paris", cfg.Timezone)
		require.NotNil(t, cfg.Author)
		require.Equal(t, "env author", *cfg.Author)
		require.Equal(t, "*", cfg.Markdown.ListStyle)
		require.True(t, cfg.Markdown.GroupsAsList)
		require.Equal(t, "issue-", cfg.Issue.DisplayPrefix)
		require.Equal(t, "^task/(\\d+)-\\S+$", cfg.Issue.Pattern)
	})

	t.Run("environment overrides project config", func(t *testing.T) {
		t.Setenv("CLG_TIMEZONE", "America/Los_Angeles")

		homeDir := t.TempDir()
		projectDir := t.TempDir()

		writeConfig(t, projectDir, "timezone: Asia/Kolkata\n")

		cfg, err := Load(homeDir, projectDir)
		require.NoError(t, err)

		require.Equal(t, "America/Los_Angeles", cfg.Timezone)
	})
}

func writeConfig(t *testing.T, dir, contents string) {
	t.Helper()

	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, maybe.WriteFile(filepath.Join(dir, configFilename), []byte(contents), 0o600))
}
