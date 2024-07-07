package niks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidID(t *testing.T) {
	require.False(t, validID(""))
	require.True(t, validID("f"))
	require.True(t, validID("foo"))
	require.True(t, validID("bar"))
	require.True(t, validID("bar123"))
	require.True(t, validID("BAR123"))
	require.False(t, validID("bar!"))
}

func TestConfig(t *testing.T) {
	var (
		p1 = Project{
			ID:   "p1",
			Name: "Hello",
		}
		p2 = Project{
			ID:   "p2",
			Name: "Hello",
		}
	)
	require.NoError(t, checkConfig(&Config{Projects: []Project{p1, p2}}))
	require.EqualError(t, checkConfig(&Config{Projects: []Project{p1, p1}}), "repeated project ID: p1")

	require.EqualError(t, checkConfig(&Config{Projects: []Project{
		Project{ID: ""},
	}}), `invalid project ID: ""`)

	require.EqualError(t, checkConfig(&Config{Projects: []Project{
		Project{ID: "foo bar"},
	}}), `invalid project ID: "foo bar"`)

}
