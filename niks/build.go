package niks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
		Start   time.Time `json:"start"`
		Done    bool      `json:"done"`
		Success bool      `json:"success"`
		Error   string    `json:"error"`
		// ExitCode int
	}
)

// this can be ordered alphabetically, and looks nice in URLs (no escaped chars)
func genID(t time.Time, projID string) string {
	return fmt.Sprintf("%s_%s", t.Format("20060102T150405"), projID)
}

// Create build ID + mkdir. This should be enough to report "in progress" in a UI.
// Use Run() after this.
func SetupBuild(buildsDir string, p Project) (*Build, error) {
	t := time.Now().UTC()
	id := genID(t, p.ID)
	path := buildsDir + id + "/"
	if err := os.MkdirAll(path, 0744); err != nil {
		return nil, err
	}

	b := &Build{
		ID:   id,
		Path: path,
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

func LoadBuild(buildsDir string, id string) (*Build, error) {
	if !validBuildDir(buildsDir, id) {
		return nil, errors.New("build not found")
	}

	return &Build{
		ID:   id,
		Path: buildsDir + id + "/",
	}, nil
}

func (b *Build) Status() Status {
	var s Status
	bytes, err := os.ReadFile(b.Path + "status.json")
	if err != nil {
		return s
	}
	json.Unmarshal(bytes, &s)
	return s
}

func (b *Build) Stdout() string {
	bytes, _ := os.ReadFile(b.Path + "stdout.txt")
	return string(bytes)
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

// Run a build. An error will be returned on a "system" error: no git available or no network. "config" errors ("no access to repo") and build errors will be part of Result. Otherwise you can read the status.json file.
// Dir structure:
//
//	.../work/ -> checkout and PWD
//	  ./stdout.txt
//	  ./status.json
func (b *Build) Run(p Project, branch string) error {
	s := Status{
		Start: time.Now().UTC(),
	}
	if err := b.WriteStatus(s); err != nil {
		return err
	}
	defer func() {
		s.Done = true
		b.WriteStatus(s)
	}()

	work := b.Path + "work/"

	stdout, err := os.Create(b.Path + "stdout.txt")
	if err != nil {
		return err
	}
	defer stdout.Close()

	call := func(cmd string, args ...string) (string, error) {
		stderr := &bytes.Buffer{}
		exe := exec.Command(cmd, args...)
		exe.Stderr = stderr
		exe.Dir = work
		stdout.WriteString("$ " + exe.String() + "\n")
		out, err := exe.Output()
		if err != nil {
			stdout.Write(stderr.Bytes())
			s.Error = fmt.Sprintf("%s: %s", cmd, err.Error())
		}
		stdout.Write(out)
		return string(out), err
	}

	{
		exe := exec.Command(cmdGit, "clone", "--depth", "1", "--branch", branch, p.Git, work)
		exe.Stdout = stdout
		exe.Stderr = stdout
		stdout.WriteString("$ " + exe.String() + "\n")
		err := exe.Run()
		if err != nil {
			s.Error = fmt.Sprintf("%s: %s", cmdGit, err.Error())
			return nil
		}
	}

	fullRev, err := call(cmdGit, "rev-parse", "HEAD")
	if err != nil {
		return nil
	}
	shortRev, err := call(cmdGit, "rev-parse", "--short", "HEAD")
	if err != nil {
		return nil
	}

	{
		exe := exec.Command(cmdNixBuild, "-A", p.Attribute)
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
			"--keep", "BRANCH_NAME",
			"--keep", "SHA",
			"--keep", "SHORT_SHA",
			"--run", p.Post,
		)
		exe := exec.Command(cmdNixShell, args...)
		exe.Stdout = stdout
		exe.Stderr = stdout
		exe.Dir = work
		exe.Env = []string{
			// fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
			fmt.Sprintf("NIX_PATH=%s", os.Getenv("NIX_PATH")),
			fmt.Sprintf("BRANCH_NAME=%s", branch),
			fmt.Sprintf("SHA=%s", fullRev), // CHECKME: don't know what the normal name is
			fmt.Sprintf("SHORT_SHA=%s", shortRev),
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
