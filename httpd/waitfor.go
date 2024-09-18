package httpd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alicebob/niksnut/niks"
)

type waitforArgs struct {
	BuildIDs []string
}

type WaitMsg struct {
	BuildID  string `json:"buildId"`
	ShortRev string `json:"shortRev"`
	Done     bool   `json:"done"`
	Success  bool   `json:"success"`
	Duration string `json:"duration"`
}

// /waitfor?buildid=123&buildid=456 -> returns a "text/event-stream" document with only the build result
func (s *Server) handlerWaitfor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.ParseForm()
	args := waitforArgs{
		BuildIDs: r.Form["buildid"],
	}
	waits, err := s.waitfor(ctx, r, &args)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Cache-Control", "no-cache")

	// This is the very simple "server-sent events" protocol.
loop:
	for {
		select {
		case st, ok := <-waits:
			if !ok {
				break loop
			}
			bits, _ := json.Marshal(st) // should have no newlines
			fmt.Fprintf(w, "data: %s\n\n", bits)
			w.(http.Flusher).Flush()
		case <-ctx.Done():
			break loop
		case <-time.After(time.Second):
			fmt.Fprintf(w, ": ping\n\n") // comment, as a keep-alive
			w.(http.Flusher).Flush()
		}
	}
	fmt.Fprintf(w, "event: finished\ndata: done\n\n")
}

func (s *Server) waitfor(ctx context.Context, r *http.Request, args *waitforArgs) (<-chan WaitMsg, error) {
	res := make(chan WaitMsg)

	builds := map[string]*niks.Build{}
	for _, id := range args.BuildIDs {
		build, err := niks.LoadBuild(s.BuildsDir, id)
		if err != nil {
			return res, err
		}
		st, _ := build.Status()
		if st.Done {
			continue
		}
		builds[id] = build
	}

	go func() {
		defer close(res)
		for {
			for id, b := range builds {
				s, _ := b.Status()
				res <- WaitMsg{
					BuildID:  id,
					ShortRev: s.ShortRev,
					Done:     s.Done,
					Success:  s.Success,
					Duration: duration(s.Finish.Sub(s.Start)),
				}
				if s.Done {
					delete(builds, id)
				}
			}
			if len(builds) == 0 {
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
		}
	}()

	return res, nil
}
