package httpd

import (
	"context"
	"net/http/httptest"
	"testing"
	"testing/fstest"

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
			Static: fstest.MapFS{
				"static/s.css": &fstest.MapFile{Data: []byte(`body { background-color: "fake"}`)},
			},
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
		m = s.Mux()
	)

	require.NoError(t, templatecheck.CheckHTML(s.loadTemplate(ctx, "index.tmpl"), indexArgs{}))
	require.NoError(t, templatecheck.CheckHTML(s.loadTemplate(ctx, "build.tmpl"), buildArgs{}))
	require.NoError(t, templatecheck.CheckHTML(s.loadTemplate(ctx, "builds.tmpl"), buildsArgs{}))

	t.Run("index", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)
		require.Equal(t, 200, res.Code)
		require.Contains(t, res.Body.String(), "proj_world")
		require.Contains(t, res.Body.String(), "proj_moon")
	})

	t.Run("static", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/static/s.css", nil)
		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)
		require.Equal(t, 200, res.Code)
	})
}
