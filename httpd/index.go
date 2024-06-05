package httpd

import (
	"net/http"

	"github.com/alicebob/niksnut/niks"
)

type indexArgs struct {
	Config niks.Config
}

func (s *Server) hndIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	args := indexArgs{
		Config: s.Config,
	}
	render(w, s.loadTemplate(ctx, "index.tmpl"), args)
}
