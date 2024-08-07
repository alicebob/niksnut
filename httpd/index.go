package httpd

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/alicebob/niksnut/niks"
)

type indexArgs struct {
	Builds       []niks.Build
	TodayYearDay int // used to show day headers in the build list
}

func (s *Server) handlerIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	args := indexArgs{
		TodayYearDay: time.Now().UTC().YearDay(),
	}
	if err := s.index(ctx, &args); err != nil {
		slog.Error("index", "error", err)
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
