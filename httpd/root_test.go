package httpd

import (
	"context"
	"testing"

	"github.com/jba/templatecheck"
	"github.com/stretchr/testify/require"
)

func TestTemplates(t *testing.T) {
	var (
		ctx = context.Background()
		s   = &Server{
			Root: staticRoot,
		}
	)
	require.NoError(t, templatecheck.CheckHTML(s.loadTemplate(ctx, "index.tmpl"), indexArgs{}))
}
