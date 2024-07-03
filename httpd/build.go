package httpd

import (
	"cmp"
	"context"
	"log/slog"
	"net/http"
	"slices"
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

// /build?projid=hello
// starts a new build
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
	args.Branch = r.FormValue("branch")
	args.Page = r.FormValue("page")

	switch args.Page {
	case "build":
		build, err := niks.SetupBuild(s.BuildsDir, args.Project)
		if err != nil {
			return err
		}
		args.BuildID = build.ID
		go func() {
			// new context, we don't want to be tied to the HTTP request
			ctx := niks.SetOffline(context.Background(), s.Offline) // FIXME: timeout
			// fixme: waitgroup
			if err := build.Run(ctx, s.BuildsDir, args.Project, args.Branch); err != nil {
				slog.Error("build failed",
					"project", args.Project,
					"branch", args.Branch,
					"error", err,
				)
			}
		}()
	case "":
		br, err := niks.Branches(ctx, s.BuildsDir, args.Project.Git)
		if err != nil {
			return err
		}
		mainFirst(br)
		args.Branches = br
		if len(br) > 0 {
			// set the default or check selection
			if args.Branch == "" || !slices.Contains(args.Branches, args.Branch) {
				args.Branch = br[0]
			}
		}
	}
	return nil
}

// Order branches, with "main" or "master" first.
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
