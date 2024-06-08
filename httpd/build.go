package httpd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alicebob/niksnut/niks"
)

type buildArgs struct {
	Page    string
	Error   string
	ProjID  string
	Project niks.Project
	Branch  string
	BuildID string
}

func (s *Server) hndBuild(w http.ResponseWriter, r *http.Request) {
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

	if args.Page == "build" {
		build, err := niks.SetupBuild(s.BuildsDir, args.Project)
		if err != nil {
			return nil
		}
		args.BuildID = build.ID
		go func() {
			// fixme: waitgroup
			// fixme: logs
			if err := build.Run(args.Project, args.Branch); err != nil {
				fmt.Printf("build: %s\n", err.Error())
			}
		}()
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
