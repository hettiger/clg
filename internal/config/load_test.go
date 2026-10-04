package config

import (
	"os"
	"path/filepath"
	"testing"

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
	require.NoError(t, os.WriteFile(filepath.Join(dir, configFilename), []byte(contents), 0o600))
}
