package httpd

import (
	"context"
	"net/http"

	"github.com/alicebob/niksnut/niks"
)

type stdoutArgs struct {
	BuildID string
	Build   niks.Build
}

func (s *Server) handlerStdout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	args := stdoutArgs{}
	if err := s.stdoud(ctx, r, &args); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(args.Build.Stdout()))
}

func (s *Server) stdoud(ctx context.Context, r *http.Request, args *stdoutArgs) error {
	id := r.FormValue("buildid")
	build, err := niks.LoadBuild(s.BuildsDir, id)
	if err != nil {
		return err
	}
	args.BuildID = id
	args.Build = *build
	return nil
}
