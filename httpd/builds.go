package httpd

import (
	"context"
	"net/http"

	"github.com/alicebob/niksnut/niks"
)

type buildsArgs struct {
	Error string
	// Project niks.Project
	// Branch  string
	BuildID string
	Build   niks.Build
	Status  niks.Status
	Stdout  string
}

func (s *Server) hndBuilds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	args := buildsArgs{}
	if err := s.builds(ctx, r, &args); err != nil {
		args.Error = err.Error()
	}
	render(w, s.loadTemplate(ctx, "builds.tmpl"), args)
}

func (s *Server) builds(ctx context.Context, r *http.Request, args *buildsArgs) error {
	id := r.FormValue("buildid")
	build, err := niks.LoadBuild(s.BuildsDir, id)
	if err != nil {
		return err
	}
	args.BuildID = id
	args.Build = *build
	args.Status = build.Status()
	args.Stdout = build.Stdout()
	return nil
}
