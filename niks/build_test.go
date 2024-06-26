package niks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsecureBuildPath(t *testing.T) {
	root := "../testdata"

	t.Run("fine", func(t *testing.T) {
		p, err := insecureBuildPath(root, "20240625T202158_hello")
		require.NoError(t, err)
		require.Equal(t, "../testdata/runs/20240625T202158_hello/", p)
	})

	t.Run("not fine", func(t *testing.T) {
		p, err := insecureBuildPath(root, "../../../../../../../etc/passwd")
		require.EqualError(t, err, "invalid build id")
		require.Equal(t, "", p)
	})
}
