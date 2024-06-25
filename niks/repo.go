package niks

import (
	"os"
)

func updateLocal(bare, repoURL string) error {
	// dir doesn't exist, or other problems
	if _, err := os.Stat(bare); err != nil {
		if err := GitCloneBare(bare, repoURL); err != nil {
			return err
		}
	}

	return GitRemoteUpdate(bare)
}

// does whatever is needed of: git clone+git fetch+git branch
func Branches(root, repoURL string) ([]string, error) {
	id := repoID(repoURL)
	bare := barePath(root, id)

	if err := updateLocal(bare, repoURL); err != nil {
		return nil, err
	}

	return GitBranches(bare)
}

// does checkout and/or fetch
func Checkout(root, repoURL, dest, branch string) error {
	id := repoID(repoURL)
	bare := barePath(root, id)

	if err := updateLocal(bare, repoURL); err != nil {
		return err
	}

	// TODO: return stderr somehow
	return GitCloneLocal(bare, dest, branch)
}
