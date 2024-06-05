package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/alicebob/niksnut/niks"
)

type (
	Result struct {
		Success bool
		Error   string
		// ExitCode int
		Stdout string
		// Stderr   string
	}
)

var buildsDir = "./builds/" // obvs should be a var
var git = "git"             //  path, fixme

// Run a build. An error will be returned on a "system" error: no git available or no network. "config" errors ("no access to repo") and build errors will be part of Result.
func Run(p niks.Project) (*Result, error) {
	t := time.Now().UTC()
	branch := "main"
	stdout := &bytes.Buffer{}

	path := fmt.Sprintf("%s%s_%s/", buildsDir, p.ID, t.Format(time.RFC3339))
	fmt.Fprintf(stdout, "path: %s\n", path)

	if err := os.MkdirAll(path, 0744); err != nil {
		return nil, err
	}

	{
		clone := exec.Command(git, "clone", "--depth", "1", "--branch", branch, p.Git, path)
		out, err := clone.Output()
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(stdout, "$ "+clone.String()+"\n"+string(out)+"\n")
	}

	{
		build := exec.Command("nix-build", "-A", p.Attribute)
		fmt.Fprintf(stdout, "$ "+build.String()+"\n")
		build.Dir = path
		build.Stdout = stdout
		build.Stderr = stdout
		if err := build.Run(); err != nil {
			return &Result{
				Stdout:  stdout.String(),
				Error:   err.Error(),
				Success: false,
			}, nil
		}
	}

	if p.Post != "" {
		post := exec.Command("bash", "-c", p.Post)
		fmt.Fprintf(stdout, "$ "+post.String()+"\n")
		post.Dir = path
		post.Stdout = stdout
		post.Stderr = stdout
		if err := post.Run(); err != nil {
			return &Result{
				Stdout:  stdout.String(),
				Error:   err.Error(),
				Success: false,
			}, nil
		}
	}

	return &Result{
		Stdout:  stdout.String(),
		Success: true,
	}, nil
}
