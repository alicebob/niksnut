package niks

import (
	"errors"
	"fmt"
	"os"
	"strings"
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
			require.EqualError(t, err, "git: could not read Username for 'https://github.com': terminal prompts disabled")
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

	t.Run("revision", func(t *testing.T) {
		full, short, err := GitRev(dest)
		require.NoError(t, err)
		require.True(t, strings.HasPrefix(full, short))
		require.NotEqual(t, full, short)
	})

	t.Run("branches", func(t *testing.T) {
		bs, err := GitBranches(dest)
		require.NoError(t, err)
		require.Contains(t, bs, "main")
		require.Contains(t, bs, "myfirstbranch")
		require.NotContains(t, bs, "")
	})
}

func TestGitError(t *testing.T) {
	err := errors.New("some error")

	// git clone https://exampleeeeee.com /tmp/bzz
	require.EqualError(t,
		gitError(err,
			`Cloning into '/tmp/bzz'...
fatal: unable to access 'https://exampleeeeee.com/': Could not resolve host: exampleeeeee.com
`,
		),
		"git: unable to access 'https://exampleeeeee.com/': Could not resolve host: exampleeeeee.com",
	)

	// git clone weirdproto://example.com /tmp/bzz
	require.EqualError(t,
		gitError(err,
			`Cloning into '/tmp/bzz'...
git: 'remote-weirdproto' is not a git command. See 'git --help'.
`,
		),
		"git: 'remote-weirdproto' is not a git command. See 'git --help'.",
	)

	// git clone weirdproto://example.com /tmp/

	require.EqualError(t,
		gitError(err,
			`fatal: destination path '/tmp' already exists and is not an empty directory.
`,
		),
		"git: destination path '/tmp' already exists and is not an empty directory.",
	)
}
