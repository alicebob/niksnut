package httpd

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"net/http"

	"github.com/alicebob/niksnut/niks"
)

type Server struct {
	Addr      string
	BuildsDir string
	Root      fs.FS // templates
	Static    fs.FS
	Config    niks.Config
	Offline   bool
}

func (s *Server) Run() error {
	h := &http.Server{
		Addr:        s.Addr,
		Handler:     s.Mux(),
		ConnContext: s.setCtx,
	}
	err := h.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Mux() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("GET /{$}", s.handlerIndex)
	m.HandleFunc("GET /build", s.handlerBuild)
	m.HandleFunc("POST /build", s.handlerBuild)
	m.HandleFunc("GET /builds", s.handlerBuilds)
	m.HandleFunc("GET /stdout", s.handlerStdout)
	m.HandleFunc("GET /stream", s.handlerStream)   // event stream of single build
	m.HandleFunc("GET /waitfor", s.handlerWaitfor) // event stream of multiple builds

	st, _ := fs.Sub(s.Static, "static")
	m.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(st)))

	return m
}

func cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=9999, immutable")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) setCtx(ctx context.Context, _ net.Conn) context.Context {
	return niks.SetOffline(ctx, s.Offline)
}

// helper to load project by an ID in handlers
func (s *Server) loadProject(id string, to *niks.Project) error {
	for _, p := range s.Config.Projects {
		if p.ID == id {
			*to = p
			return nil
		}
	}
	return errors.New("no such project")
}
