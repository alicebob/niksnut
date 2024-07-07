package httpd

import (
	"embed"
	"net/http"
	"os"
)

var (
	//go:embed *.tmpl
	TemplateRoot embed.FS
)

func httpUser(r *http.Request) string {
	if user := r.Header.Get("X-Remote-User"); user != "" {
		return user
	}
	// for tests, mostly
	return os.Getenv("NIKSNUTUSER")
}
