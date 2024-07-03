package niks

import (
	"context"
	"errors"
	"os"
)

// updateLocal does a clone if needed, and then a fetch.
// In offline mode this will fail if it needs to do a git clone.
func updateLocal(ctx context.Context, bare, repoURL string) error {
	// dir doesn't exist, or other problems
	if _, err := os.Stat(bare); err != nil {
		if isOffline(ctx) {
			return errors.New("cannot run git clone in offline mode")
		}
		if err := GitCloneBare(ctx, bare, repoURL); err != nil {
			return err
		}
	}

	if isOffline(ctx) {
		return nil
	}

	return GitRemoteUpdate(ctx, bare)
}

// does whatever is needed of: git clone+git fetch+git branch
func Branches(ctx context.Context, root, repoURL string) ([]string, error) {
	id := repoID(repoURL)
	bare := barePath(root, id)

	if err := updateLocal(ctx, bare, repoURL); err != nil {
		return nil, err
	}

	return GitBranches(ctx, bare)
}

// does checkout and/or fetch
func Checkout(ctx context.Context, root, repoURL, dest, branch string) error {
	id := repoID(repoURL)
	bare := barePath(root, id)

	if err := updateLocal(ctx, bare, repoURL); err != nil {
		return err
	}

	return GitCloneLocal(ctx, bare, dest, branch)
}
