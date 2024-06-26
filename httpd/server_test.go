package httpd

import (
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	"github.com/alicebob/niksnut/niks"
)

func TestMux(t *testing.T) {
	var (
		s = &Server{
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

	t.Run("index", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)
		require.Equal(t, 200, res.Code)
		require.Contains(t, res.Body.String(), "proj_world")
		require.Contains(t, res.Body.String(), "proj_moon")
	})

	t.Run("builds", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/builds?buildid=20240625T202158_hello", nil)
		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)
		require.Equal(t, 200, res.Code)
		require.Contains(t, res.Body.String(), "Success")
		require.Contains(t, res.Body.String(), "Stdout")
	})

	t.Run("stdout", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/stdout?buildid=20240625T202158_hello", nil)
		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)
		require.Equal(t, 200, res.Code)
		require.Contains(t, res.Body.String(), "ssh version")
	})

	t.Run("static", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/static/s.css", nil)
		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)
		require.Equal(t, 200, res.Code)
	})
}
