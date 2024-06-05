package httpd

import (
	"net/http"
)

type indexArgs struct {
}

func (s *Server) hndIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	render(w, s.loadTemplate(ctx, "index.tmpl"), nil)
}
