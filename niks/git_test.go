package niks

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepoID(t *testing.T) {
	require.Equal(t, "ssh_github.com_alicebob_gohello", repoID("ssh://git@github.com/alicebob/gohello"))
	require.Equal(t, "https_github.com_alicebob123_gohello", repoID("https://github.com/alicebob123/gohello"))
}

func TestGit(t *testing.T) {
	root := SetupInt(t)

	t.Run("clone remote", func(t *testing.T) {
		t.Run("no such repo", func(t *testing.T) {
			repo := "https://github.com/alicebob/gohellonosuch"
			dest := barePath(root, repo)
			err := GitCloneBare(dest, repo)
			require.EqualError(t, err, "exit status 128")
		})

		t.Run("fine", func(t *testing.T) {
			repo := "https://github.com/alicebob/gohello"
			dest := root + "/clone"
			err := GitCloneBare(dest, repo)
			require.NoError(t, err)

			// quick check if the repo is there
			_, err = os.Stat(fmt.Sprintf("%s/config", dest))
			require.NoError(t, err)
		})
	})

	repo := "https://github.com/alicebob/gohello"
	dest := barePath(root, repo)
	require.NoError(t, GitCloneBare(dest, repo))

	t.Run("clone local", func(t *testing.T) {
		t.Run("fine", func(t *testing.T) {
			local := fmt.Sprintf("%s/clonetest", root)
			require.NoError(t, GitCloneLocal(dest, local, "main"))

			_, err := os.Stat(fmt.Sprintf("%s/.git/config", local))
			require.NoError(t, err)
			_, err = os.Stat(fmt.Sprintf("%s/README.md", local))
			require.NoError(t, err)
		})
	})

	t.Run("remote", func(t *testing.T) {
		require.NoError(t, GitRemoteUpdate(dest))
	})

	t.Run("branches", func(t *testing.T) {
		bs, err := GitBranches(dest)
		require.NoError(t, err)
		require.Contains(t, bs, "main")
		require.Contains(t, bs, "myfirstbranch")
		require.NotContains(t, bs, "")
	})
}
