package niks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBranches(t *testing.T) {
	root := SetupInt(t)
	repoURL := "https://github.com/alicebob/gohello"

	{
		branches, err := Branches(root, repoURL)
		require.NoError(t, err)
		require.Contains(t, branches, "main")
		require.NotContains(t, branches, "")
	}
	// again is fine
	{
		branches, err := Branches(root, repoURL)
		require.NoError(t, err)
		require.Contains(t, branches, "main")
		require.NotContains(t, branches, "")
	}
}

func TestCheckout(t *testing.T) {
	root := SetupInt(t)
	repoURL := "https://github.com/alicebob/gohello"
	require.NoError(t, Checkout(root, repoURL, root+"/checkedout_main", "main"))
	require.NoError(t, Checkout(root, repoURL, root+"/checkedout_myfirstbranch", "myfirstbranch"))
}
