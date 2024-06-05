package httpd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alicebob/niksnut/niks"
)

type indexArgs struct {
	Config niks.Config
	Builds []niks.Build
}

func (s *Server) hndIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	args := indexArgs{
		Config: s.Config,
	}
	if err := s.index(ctx, &args); err != nil {
		fmt.Printf("index: %s\n", err.Error()) // FIXME slog
	}
	render(w, s.loadTemplate(ctx, "index.tmpl"), args)
}

func (s *Server) index(ctx context.Context, args *indexArgs) error {
	ls, err := niks.ListBuilds(s.BuildsDir)
	if err != nil {
		return err
	}
	args.Builds = ls
	return nil
}
