package httpd

import (
	"embed"
)

var (
	//go:embed *.tmpl
	staticRoot embed.FS
)
