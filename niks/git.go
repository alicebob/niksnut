package niks

// git helpers, which use the `git` tool.
// Some of these could be implemented with https://pkg.go.dev/github.com/go-git/go-git/v5#example-Clone but that's such a big tree of dependencies.

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// makes a "path safe id" from a repo. It's not a 100% unique id, but Good Enough.
func repoID(repo string) string {
	u, err := url.Parse(repo)
	if err != nil {
		slog.Error("invalid repo id, this will cause errors", "repo", repo)
		return "invalidrepo"
	}

	id := (&url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	}).String()
	id = strings.ToLower(id)
	return regexp.MustCompile(`[^a-z0-9.]+`).ReplaceAllString(id, "_")
	return id
}

// path to bare repo
func barePath(root, repoURL string) string {
	id := repoID(repoURL)
	return fmt.Sprintf("%s/repo/%s.git", root, id)
}

// dest should come from RepoPath()
func GitCloneBare(dest, repoURL string) error {
	exe := exec.Command(cmdGit, "clone", "--bare", "--mirror", repoURL, dest)
	exe.Env = []string{
		"GIT_TERMINAL_PROMPT=0",
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	}
	stdout, err := exe.CombinedOutput()
	if err != nil {
		slog.Error("git clone",
			"cmd", exe.String(),
			"stdout", stdout,
			"error", err,
		)
		return err
	}
	return nil
}

func GitRemoteUpdate(repo string) error {
	exe := exec.Command(cmdGit, "remote", "update", "--prune")
	exe.Dir = repo
	exe.Env = []string{
		"GIT_TERMINAL_PROMPT=0",
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	}
	stdout, err := exe.CombinedOutput()
	if err != nil {
		slog.Error("git remote",
			"cmd", exe.String(),
			"stdout", stdout,
			"error", err,
		)
		return err
	}
	return nil
}

// local clone of a bare clone
func GitCloneLocal(src, dest, branch string) error {
	exe := exec.Command(cmdGit, "clone", "--branch", branch, src, dest)
	stdout, err := exe.CombinedOutput()
	if err != nil {
		slog.Error("git clone",
			"cmd", exe.String(),
			"stdout", stdout,
			"error", err,
		)
		return err
	}
	return nil
}

func GitBranches(src string) ([]string, error) {
	exe := exec.Command(cmdGit, "for-each-ref", "--format", "%(refname:short)", "refs/heads/")
	exe.Dir = src
	stdout, err := exe.CombinedOutput()
	if err != nil {
		slog.Error("git branch",
			"cmd", exe.String(),
			"stdout", stdout,
			"error", err,
		)
		return nil, err
	}
	var bs []string
	for _, b := range strings.Split(string(stdout), "\n") {
		br := strings.TrimSpace(b)
		if br != "" {
			bs = append(bs, br)
		}
	}
	return bs, nil
}
