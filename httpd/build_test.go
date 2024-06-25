package httpd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMainFirst(t *testing.T) {
	s := []string{"q", "main", "a", "c", "d", "aa", "Aa", "Bb"}
	mainFirst(s)
	require.Equal(t,
		[]string{"main", "Aa", "Bb", "a", "aa", "c", "d", "q"},
		s,
	)
}
