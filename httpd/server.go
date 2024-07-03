package httpd

import (
	"errors"
	"io/fs"
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
		Addr:    s.Addr,
		Handler: s.Mux(),
	}
	err := h.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Mux() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("GET /{$}", s.setCtx(s.handlerIndex))
	m.HandleFunc("GET /build", s.setCtx(s.handlerBuild))
	m.HandleFunc("POST /build", s.setCtx(s.handlerBuild))
	m.HandleFunc("GET /builds", s.setCtx(s.handlerBuilds))
	m.HandleFunc("GET /stdout", s.handlerStdout)
	m.HandleFunc("GET /stream", s.handlerStream)

	st, _ := fs.Sub(s.Static, "static")
	m.Handle("GET /static/*", http.StripPrefix("/static/", cache(http.FileServerFS(st))))

	return m
}

func cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=9999, immutable")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) setCtx(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := niks.SetOffline(r.Context(), s.Offline)
		r = r.WithContext(ctx)
		next(w, r)
	})
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
