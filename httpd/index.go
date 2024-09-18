package httpd

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/alicebob/niksnut/niks"
)

type indexArgs struct {
	TodayYearDay int // used to show day headers in the build list
	Builds       []niks.Build
	Building     []string // we can ./waitfor these
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

	for _, b := range ls {
		if s, _ := b.Status(); !s.Done {
			// we could skip old IDs which hang for whatever reason
			args.Building = append(args.Building, b.ID)
		}
	}

	return nil
}
