package httpd

import (
	"errors"
	"io/fs"
	"net/http"
)

type Server struct {
	Addr string
	Root fs.FS
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
	m.HandleFunc("GET /", s.hndIndex)
	return m
}
