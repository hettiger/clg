package changelog_test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestGitService_CurrentBranch(t *testing.T) {
	tests := []struct {
		name    string
		branch  string
		want    string
		wantErr bool
	}{
		{
			name:   "existing repo",
			branch: "fake-branch",
			want:   "fake-branch",
		},
		{
			name:    "missing repo",
			branch:  "",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cwd := t.TempDir()
			g := changelog.NewGitService(cwd)

			if tt.branch != "" {
				cmd := exec.Command("git", "init", "-b", tt.branch)
				cmd.Dir = cwd
				err := cmd.Run()
				require.NoError(t, err)
			}

			got, gotErr := g.CurrentBranch()

			if tt.wantErr {
				require.Error(t, gotErr)
			} else {
				require.NoError(t, gotErr)
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestGitService_AuthorName(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		want    string
		wantErr bool
	}{
		{
			name:    "author",
			config:  filepath.Join("testdata", "gitconfig_author"),
			want:    "Fake Author",
			wantErr: false,
		},
		{
			name:    "empty author",
			config:  filepath.Join("testdata", "gitconfig_empty_author"),
			want:    "",
			wantErr: false,
		},
		{
			name:    "missing author",
			config:  filepath.Join("testdata", "gitconfig_missing_author"),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath, err := filepath.Abs(tt.config)
			require.NoError(t, err)
			t.Setenv("GIT_CONFIG_GLOBAL", configPath)
			t.Setenv("GIT_CONFIG_NOSYSTEM", "1")

			d := t.TempDir()
			g := changelog.NewGitService(d)

			got, gotErr := g.AuthorName()

			if tt.wantErr {
				require.Error(t, gotErr)
			} else {
				require.NoError(t, gotErr)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
