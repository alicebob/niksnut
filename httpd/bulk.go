package httpd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/alicebob/niksnut/niks"
)

type bulkArgs struct {
	Error    string
	Projects []niks.Project
	Selected map[string]bool // project ID -> selected
	Branches []string        // branches common to all selected projects
	Branch   string
	BuildIDs []string // IDs of builds started
}

// /bulk - show bulk deployment page
func (s *Server) handlerBulk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	args := bulkArgs{
		Projects: s.Config.Projects,
		Selected: map[string]bool{},
	}
	if err := s.bulk(ctx, r, &args); err != nil {
		args.Error = err.Error()
	}
	render(w, s.loadTemplate(ctx, "bulk.tmpl"), args)
}

func (s *Server) bulk(ctx context.Context, r *http.Request, args *bulkArgs) error {
	// Handle form submission
	if r.Method == "POST" {
		for _, proj := range args.Projects {
			if r.FormValue("project-"+proj.ID) == "on" {
				args.Selected[proj.ID] = true
			}
		}

		switch r.FormValue("page") {
		case "deploy":
			args.Branch = r.FormValue("branch")
			if args.Branch == "" {
				return nil
			}

			s.startBulkBuilds(ctx, args)
			return nil
		case "select":
			// Only fetch branches when projects are actually selected
			branchesByRepo := map[string][]string{}
			for _, proj := range args.Projects {
				if _, ok := args.Selected[proj.ID]; ok {
					if _, ok := branchesByRepo[proj.Git]; ok {
						continue
					}

					br, err := niks.Branches(ctx, s.BuildsDir, proj.Git)
					if err != nil {
						return fmt.Errorf("error getting branches for project %s: %w", proj.ID, err)
					}
					branchesByRepo[proj.Git] = br
				}
			}
			args.Branches = computeCommonBranches(args.Projects, args.Selected, branchesByRepo)
		}
	}

	return nil
}

// computeCommonBranches returns branches common to all selected projects
// when branches are stored by repo URL instead of project ID
func computeCommonBranches(projects []niks.Project, selected map[string]bool, branches map[string][]string) []string {
	if len(selected) == 0 {
		return nil
	}

	// Collect unique repo URLs for selected projects
	repoURLs := map[string]bool{}
	for _, proj := range projects {
		if selected[proj.ID] {
			repoURLs[proj.Git] = true
		}
	}

	var common []string
	for repoURL := range repoURLs {
		if br, ok := branches[repoURL]; ok {
			if common == nil {
				// First repo - use all its branches
				common = br
			} else {
				// Subsequent repos - intersect with existing common branches
				common = slices.DeleteFunc(common, func(branch string) bool {
					return !slices.Contains(br, branch)
				})
				if len(common) == 0 {
					break
				}
			}
		}
	}
	return common
}

func (s *Server) startBulkBuilds(ctx context.Context, args *bulkArgs) {
	var ids []string
	for projID, on := range args.Selected {
		if !on {
			continue
		}
		var proj niks.Project
		if err := s.loadProject(projID, &proj); err != nil {
			slog.Error("project not found", "id", projID, "error", err)
			continue
		}

		build, err := niks.SetupBuild(s.BuildsDir, proj)
		if err != nil {
			slog.Error("error setting up build", "project", proj.ID, "error", err)
			continue
		}
		ids = append(ids, build.ID)

		go func(p niks.Project, b *niks.Build, branch string) {
			// new context, we don't want to be tied to the HTTP request
			buildCtx := niks.SetOffline(context.Background(), s.Offline)
			if err := b.Run(buildCtx, s.BuildsDir, p, branch, "bulk-user"); err != nil {
				slog.Error("build failed",
					"project", p.ID,
					"branch", branch,
					"error", err,
				)
			}
		}(proj, build, args.Branch)
	}
	args.BuildIDs = ids
}
