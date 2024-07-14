package niks

import (
	"context"
	"log/slog"
	"time"
)

func GarbageCollect(ctx context.Context, root string, keepAfter time.Time) error {
	bs, err := ListBuilds(root)
	if err != nil {
		return err
	}

	for _, b := range bs {
		st, err := b.Status()
		if err != nil {
			slog.Info("gc build load failed, deleting", "buildid", b.ID, "error", err)
			if err := RemoveBuild(root, b.ID); err != nil {
				return err
			}
			continue
		}
		if st.Start.Before(keepAfter) {
			slog.Info("gc build cleanup", "buildid", b.ID)
			if err := RemoveBuild(root, b.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
