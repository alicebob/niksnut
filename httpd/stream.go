package httpd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alicebob/niksnut/niks"
)

type streamArgs struct {
	BuildID string
	Build   niks.Build
	Offset  int
}

// /stream?buildid=123&offset=123 -> returns a "text/event-stream" document.
func (s *Server) handlerStream(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	args := streamArgs{}
	if err := s.stream(ctx, r, &args); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	slog.Info("start stream")
	defer slog.Info("stop stream")

	w.Header().Set("Content-Type", "text/event-stream")
	// MDN docs set these two, and who am I to question that?
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Cache-Control", "no-cache")

	// The "server-sent events" protocol is a stream of "data: foobar"
	// messages, with an empty line after. We can do that.
	offset := args.Offset
loop:
	for {
		fmt.Fprintf(w, ": ping\n") // comment, as a keep-alive
		bs := args.Build.StdoutOffset(offset)
	lines:
		for {
			// only show complete lines.
			car, cdr, found := strings.Cut(bs, "\n")
			if !found {
				break lines
			}
			fmt.Fprintf(w, "data: %s\n", car)
			offset += len(car) + 1
			bs = cdr
		}
		fmt.Fprintf(w, "\n")
		w.(http.Flusher).Flush()

		if s, _ := args.Build.Status(); s.Done {
			break loop
		}

		select {
		case <-ctx.Done():
			break loop
		case <-time.After(time.Second):
		}
	}

	// We end with a custom message with the final status, so the UI could update
	// that (it simple reloads when it gets the 'finished' event).
	fmt.Fprintf(w, "event: finished\n")
	st, _ := args.Build.Status()
	fmt.Fprintf(w, `data: {"done": %t, "success": %t}`, st.Done, st.Success)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\n")
}

func (s *Server) stream(ctx context.Context, r *http.Request, args *streamArgs) error {
	id := r.FormValue("buildid")
	build, err := niks.LoadBuild(s.BuildsDir, id)
	if err != nil {
		return err
	}
	args.BuildID = id
	args.Build = *build
	if o := r.FormValue("offset"); o != "" {
		n, _ := strconv.Atoi(o)
		args.Offset = n
	}

	return nil
}
