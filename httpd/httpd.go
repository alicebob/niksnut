package httpd

import (
	"embed"
)

var (
	//go:embed *.tmpl
	TemplateRoot embed.FS
)
