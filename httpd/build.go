package httpd

import (
	"cmp"
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/alicebob/niksnut/niks"
)

type buildArgs struct {
	Page     string
	Error    string
	ProjID   string
	Project  niks.Project
	Branch   string
	BuildID  string
	Branches []string
}

func (s *Server) handlerBuild(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	args := buildArgs{}
	if err := s.build(ctx, r, &args); err != nil {
		args.Error = err.Error()
		args.Page = ""
	}
	render(w, s.loadTemplate(ctx, "build.tmpl"), args)
}

func (s *Server) build(ctx context.Context, r *http.Request, args *buildArgs) error {
	id := r.FormValue("projid")
	if err := s.loadProject(id, &args.Project); err != nil {
		return err
	}
	args.ProjID = id
	args.Branch = def(r.FormValue("branch"), "main")
	args.Page = r.FormValue("page")

	switch args.Page {
	case "build":
		build, err := niks.SetupBuild(s.BuildsDir, args.Project)
		if err != nil {
			return nil
		}
		args.BuildID = build.ID
		go func() {
			// fixme: waitgroup
			// fixme: logs
			if err := build.Run(s.BuildsDir, args.Project, args.Branch); err != nil {
				fmt.Printf("build: %s\n", err.Error())
			}
		}()
	case "":
		br, err := niks.Branches(s.BuildsDir, args.Project.Git)
		if err != nil {
			return err
		}
		mainFirst(br)
		args.Branches = br
	}
	return nil
}

func def[A comparable](v, d A) A {
	var zero A
	if v == zero {
		return d
	}
	return v
}

// Order branches, with "main" or "master" first.
// (we should get from git what that branch is)
func mainFirst(bs []string) {
	sort.Slice(bs, func(i, j int) bool {
		if bs[i] == "main" || bs[i] == "master" {
			return true
		}
		if bs[j] == "main" || bs[j] == "master" {
			return false
		}
		return cmp.Less(bs[i], bs[j])
	})
}
