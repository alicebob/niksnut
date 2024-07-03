package niks

import (
	"context"
	"os/exec"
)

func nixCollectGarbage(ctx context.Context) (string, error) {
	exe := exec.CommandContext(ctx, cmdNixGC)
	stdout, err := exe.CombinedOutput()
	return string(stdout), err
}
