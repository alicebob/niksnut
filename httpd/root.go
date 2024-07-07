package httpd

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/alicebob/niksnut/niks"
)

func (s *Server) rootTemplate() *template.Template {
	root := template.New("root").Funcs(template.FuncMap{
		"showerror":  showerror,
		"showstatus": showstatus,
		"radio":      htmlRadio,
		"datetime":   datetime,
		"duration":   duration,
		"project": func(projID string) *niks.Project {
			var p niks.Project
			if err := s.loadProject(projID, &p); err != nil {
				return nil
			}
			return &p
		},
		"static": s.staticLink,
		"config": func() niks.Config { return s.Config },
		// ...
	})
	return template.Must(root.ParseFS(s.Root, "root.tmpl"))
}

func (s *Server) loadTemplate(ctx context.Context, name string) *template.Template {
	b, err := fs.ReadFile(s.Root, name)
	if err != nil {
		panic(err)
	}
	return template.Must(s.rootTemplate().Parse(string(b)))
}

// returns link l but with ?sha1=foobar attached
// This allows very strict cacheing headers, while still always getting the latest version in the client.
func (s *Server) staticLink(l string) string {
	bits, err := fs.ReadFile(s.Static, strings.TrimPrefix(l, "/"))
	if err != nil {
		panic(err)
	}
	sum := sha1.Sum(bits)
	return l + "?sha1=" + hex.EncodeToString(sum[:])
}

func render(w http.ResponseWriter, t *template.Template, args interface{}) {
	b := &bytes.Buffer{}
	if err := t.ExecuteTemplate(b, "root", args); err != nil {
		fmt.Printf("render template: %s\n", err.Error())
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.Write(b.Bytes())
}
