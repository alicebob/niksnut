package niks

import (
	"context"
)

var (
	cmdGit            = "git"
	cmdNixInstantiate = "nix-instantiate"
	cmdNixBuild       = "nix-build"
	cmdNixShell       = "nix-shell"
	cmdNixGC          = "nix-collect-garbage"

	ctxOffline = "niksoffline"
)

// context can indicate we're in offline mode, which is nice when you're working in a train without wifi.
func SetOffline(ctx context.Context, offline bool) context.Context {
	return context.WithValue(ctx, ctxOffline, offline)
}

func isOffline(ctx context.Context) bool {
	v := ctx.Value(ctxOffline)
	if v == nil {
		return false
	}
	if t, ok := v.(bool); ok {
		return t
	}
	return false
}
