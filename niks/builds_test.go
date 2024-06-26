package niks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListBuilds(t *testing.T) {
	ls, err := ListBuilds("../testdata")
	require.NoError(t, err)
	require.Len(t, ls, 3)
}
