package niks

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type (
	Build struct {
		ID         string
		Path       string
		StdoutFile string
		StatusFile string
	}
	Status struct {
		Start   time.Time `json:"start"`
		Done    bool      `json:"done"`
		Success bool      `json:"success"`
		Error   string    `json:"error"`
		// ExitCode int
	}
)

var buildsDir = "./builds/" // obvs should be a var
var git = "git"             //  path, fixme

// Create build ID + mkdir. This should be enough to report "in progress" in a UI.
// Use Run() after this.
func SetupBuild(p Project) (*Build, error) {
	t := time.Now().UTC()
	id := fmt.Sprintf("%s_%s", t.Format(time.RFC3339), p.ID)
	path := buildsDir + id + "/"
	if err := os.MkdirAll(path, 0744); err != nil {
		return nil, err
	}

	b := &Build{
		ID:         id,
		Path:       path,
		StdoutFile: path + "stdout.txt",
		StatusFile: path + "status.json",
	}

	// create the files to make reading them easy
	if err := os.WriteFile(b.StdoutFile, nil, 0666); err != nil {
		return nil, err
	}
	if err := b.WriteStatus(Status{}); err != nil {
		return nil, err
	}
	return b, nil
}

// func (b *Build) Status() Status {
// }
func (b *Build) Stdout() string {
	bytes, _ := os.ReadFile(b.StdoutFile)
	return string(bytes)
}

func (b *Build) WriteStatus(s Status) error {
	fh, err := os.Create(b.StatusFile)
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

	stdout, err := os.Create(b.StdoutFile)
	if err != nil {
		return err
	}
	defer stdout.Close()

	{
		clone := exec.Command(git, "clone", "--depth", "1", "--branch", branch, p.Git, work)
		out, err := clone.Output()
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "$ "+clone.String()+"\n"+string(out)+"\n")
	}

	{
		build := exec.Command("nix-build", "-A", p.Attribute)
		fmt.Fprintf(stdout, "$ "+build.String()+"\n")
		build.Dir = work
		build.Stdout = stdout
		build.Stderr = stdout
		if err := build.Run(); err != nil {
			s.Error = err.Error()
			return nil
		}
	}

	if p.Post != "" {
		post := exec.Command("bash", "-c", p.Post)
		fmt.Fprintf(stdout, "$ "+post.String()+"\n")
		post.Dir = work
		post.Stdout = stdout
		post.Stderr = stdout
		if err := post.Run(); err != nil {
			s.Error = err.Error()
			return nil
		}
	}

	s.Success = true

	return nil
}
