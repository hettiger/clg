package changelog_test

import (
	"os/exec"
	"testing"

	"github.com/hettiger/clg/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestGitServcie_CurrentBranch(t *testing.T) {
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
