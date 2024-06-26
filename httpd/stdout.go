package httpd

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/alicebob/niksnut/niks"
)

type stdoutArgs struct {
	BuildID string
	Build   niks.Build
	Tail    bool
}

// /stdout?buildid=123 -> returns plain text stdout
// If ?tail=true is set this "streams" the log (simple polling). Note that browsers don't render anything until the document is a certain size, but after that it works well.
func (s *Server) handlerStdout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	args := stdoutArgs{}
	if err := s.stdoud(ctx, r, &args); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	pos := 0
	bs := args.Build.Stdout()
	n, _ := w.Write([]byte(bs))
	pos += n

	// tail -f is implemented by polling, while checking the "done" flag.
	if args.Tail {
		slog.Info("starting tail -f")
		defer slog.Info("stopped tail -f")
	tail:
		for {
			select {
			case <-ctx.Done():
				break tail
			case <-time.After(time.Second):
			}
			if args.Build.Status().Done {
				break tail
			}
			bs := args.Build.StdoutOffset(pos)
			n, _ := w.Write([]byte(bs))
			pos += n
			w.(http.Flusher).Flush()
		}
	}
}

func (s *Server) stdoud(ctx context.Context, r *http.Request, args *stdoutArgs) error {
	id := r.FormValue("buildid")
	build, err := niks.LoadBuild(s.BuildsDir, id)
	if err != nil {
		return err
	}
	args.BuildID = id
	args.Build = *build

	args.Tail = r.FormValue("tail") != ""

	return nil
}
