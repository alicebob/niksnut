package httpd

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

func (s *Server) rootTemplate() *template.Template {
	root := template.New("root").Funcs(template.FuncMap{
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
