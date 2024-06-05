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
	Root      fs.FS
	Config    niks.Config
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
	m.HandleFunc("GET /{$}", s.hndIndex)
	m.HandleFunc("GET /build", s.hndBuild)
	m.HandleFunc("POST /build", s.hndBuild)
	m.HandleFunc("GET /builds", s.hndBuilds)
	return m
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
