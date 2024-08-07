package niks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type (
	Build struct {
		ID   string
		Path string // with trailing /
	}
	Status struct {
		ProjectID string    `json:"projId"`
		User      string    `json:"user"`
		Branch    string    `json:"branch"`
		Rev       string    `json:"rev"`
		ShortRev  string    `json:"shortRev"`
		Start     time.Time `json:"start"`
		Finish    time.Time `json:"finish"`
		Done      bool      `json:"done"`
		Success   bool      `json:"success"`
		Error     string    `json:"error"`
		// ExitCode int
	}
)

// this can be ordered alphabetically, and looks nice in URLs (no escaped chars)
func genID(t time.Time, projID string) string {
	return fmt.Sprintf("%s_%s", t.Format("20060102T150405"), projID)
}

func buildPath(root, id string) string {
	return fmt.Sprintf("%s/runs/%s/", root, id)
}

// like buildPath(), but checks if "id" exists.
func insecureBuildPath(root, id string) (string, error) {
	runs := fmt.Sprintf("%s/runs/", root)
	ls, err := os.ReadDir(runs)
	if err != nil {
		return "", err
	}
	for _, l := range ls {
		if l.Name() == id {
			return fmt.Sprintf("%s%s/", runs, id), nil
		}
	}
	return "", errors.New("invalid build id")
}

// Create build ID + mkdir. This should be enough to report "in progress" in a UI.
// Use Run() after this.
func SetupBuild(root string, p Project) (*Build, error) {
	t := time.Now().UTC()
	id := genID(t, p.ID)
	dst := buildPath(root, id)
	if err := os.MkdirAll(dst, 0744); err != nil {
		return nil, err
	}

	b := &Build{
		ID:   id,
		Path: dst,
	}

	// create the files to make reading them easy
	if err := os.WriteFile(b.Path+"stdout.txt", nil, 0666); err != nil {
		return nil, err
	}
	// it's only a valid build dir with a status.json file
	if err := b.WriteStatus(Status{}); err != nil {
		return nil, err
	}
	return b, nil
}

func LoadBuild(root string, id string) (*Build, error) {
	dst, err := insecureBuildPath(root, id)
	if err != nil {
		return nil, errors.New("build not found")
	}
	if !validBuildDir(dst) {
		return nil, errors.New("build not found")
	}

	return &Build{
		ID:   id,
		Path: dst,
	}, nil
}

// Deleted the complete build dir (we don't keep a DB or anything, so that's
// all).
// Needs an untainted build ID (so it shouldn't come from the UI).
func RemoveBuild(root string, buildID string) error {
	p := buildPath(root, buildID)
	return os.RemoveAll(p)
}

func (b *Build) Status() (Status, error) {
	var s Status
	bytes, err := os.ReadFile(b.Path + "status.json")
	if err != nil {
		return s, err
	}
	return s, json.Unmarshal(bytes, &s)
}

func (b *Build) Stdout() string {
	bs, _ := os.ReadFile(b.stdoutFile())
	return string(bs)
}

func (b *Build) StdoutOffset(n int) string {
	fh, err := os.Open(b.stdoutFile())
	if err != nil {
		return ""
	}
	defer fh.Close()
	fh.Seek(int64(n), io.SeekStart)
	bs, _ := io.ReadAll(fh)
	return string(bs)
}

func (b *Build) stdoutFile() string {
	return b.Path + "stdout.txt"
}

func (b *Build) WriteStatus(s Status) error {
	fh, err := os.Create(b.Path + "status.json")
	if err != nil {
		return err
	}
	defer fh.Close()
	e := json.NewEncoder(fh)
	e.SetIndent("", "  ")
	return e.Encode(s)
}

// Run a build. An error will be returned on a "system" error: no git available or no network. "config" errors ("no access to repo") and build errors will be part of Result. For other info you can read the status.json file.
// Dir structure:
//
//	/repo/repoid.git -> headless checkout
//	/runs/123/work/ -> checkout and PWD
//	        ./stdout.txt
//	        ./status.json
func (b *Build) Run(ctx context.Context, root string, p Project, branch, user string) error {
	s := Status{
		ProjectID: p.ID,
		User:      user,
		Branch:    branch,
		Start:     time.Now().UTC(),
	}
	if err := b.WriteStatus(s); err != nil {
		return err
	}
	defer func() {
		s.Done = true
		s.Finish = time.Now().UTC()
		b.WriteStatus(s)
	}()

	work := b.Path + "work/"

	stdout, err := os.Create(b.Path + "stdout.txt")
	if err != nil {
		return err
	}
	defer stdout.Close()

	if err := Checkout(ctx, root, p.Git, work, branch); err != nil {
		s.Error = fmt.Sprintf("checkout: %s", err.Error())
		return nil // returns not an error
	}

	fullRev, shortRev, err := GitRev(ctx, work)
	if err != nil {
		return err
	}
	s.Rev = fullRev
	s.ShortRev = shortRev
	b.WriteStatus(s)

	{
		exe := exec.CommandContext(ctx, cmdNixBuild, "-A", p.Attribute)
		exe.Stdout = stdout
		exe.Stderr = stdout
		exe.Dir = work
		stdout.WriteString("$ " + exe.String() + "\n")
		err := exe.Run()
		if err != nil {
			s.Error = fmt.Sprintf("%s: %s", cmdNixBuild, err.Error())
			return nil
		}
	}

	if p.Post != "" {
		args := []string{}
		if len(p.Packages) > 0 {
			args = append(args, "-p")
			args = append(args, p.Packages...)
		}
		args = append(args, "--pure",
			"--keep", "HOME",
			"--keep", "USER",
			"--keep", "BRANCH_NAME",
			"--keep", "SHA",
			"--keep", "SHORT_SHA",
			"--run", p.Post,
		)
		exe := exec.CommandContext(ctx, cmdNixShell, args...)
		exe.Stdout = stdout
		exe.Stderr = stdout
		exe.Dir = work
		exe.Env = []string{
			// fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
			fmt.Sprintf("USER=%s", os.Getenv("USER")),
			fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
			fmt.Sprintf("BRANCH_NAME=%s", branch),
			fmt.Sprintf("SHA=%s", fullRev), // CHECKME: don't know what the normal name is
			fmt.Sprintf("SHORT_SHA=%s", shortRev),
		}
		if p := os.Getenv("NIX_PATH"); p != "" {
			exe.Env = append(exe.Env, fmt.Sprintf("NIX_PATH=%s", p))
		}
		stdout.WriteString("$ " + exe.String() + "\n")
		if err := exe.Run(); err != nil {
			s.Error = err.Error()
			return nil
		}
	}

	s.Success = true

	return nil
}
