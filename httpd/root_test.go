package httpd

import (
	"context"
	"testing"

	"github.com/jba/templatecheck"
	"github.com/stretchr/testify/require"

	"github.com/alicebob/niksnut/niks"
)

func TestTemplates(t *testing.T) {
	var (
		ctx = context.Background()
		s   = &Server{
			Root:      TemplateRoot,
			BuildsDir: "../testdata/",
			Config: niks.Config{
				Projects: []niks.Project{
					{
						ID:   "proj_world",
						Name: "Save the world",
					},
					{
						ID:   "proj_moon",
						Name: "Save the moon",
					},
				},
			},
		}
	)

	require.NoError(t, templatecheck.CheckHTML(s.loadTemplate(ctx, "index.tmpl"), indexArgs{}))
	require.NoError(t, templatecheck.CheckHTML(s.loadTemplate(ctx, "build.tmpl"), buildArgs{}))
	require.NoError(t, templatecheck.CheckHTML(s.loadTemplate(ctx, "builds.tmpl"), buildsArgs{}))
}
