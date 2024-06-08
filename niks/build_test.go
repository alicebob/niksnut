package niks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindDerivs(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		s := "\n\t\t\t\techo ENV: $(printenv)\n\t\t\t\t$(/nix/store/0g1s8yd0biawp32fl3i7kdbi219jx6aq-openssh-9.7p1/bin/ssh --version)\n\t\t\t"
		drvs := findDerivs(s)
		require.Equal(t, []string{"/nix/store/0g1s8yd0biawp32fl3i7kdbi219jx6aq-openssh-9.7p1"}, drvs)
	})
	t.Run("space", func(t *testing.T) {
		s := "hello /nix/store/zhpnpc6ljvdp10ky2wfqwph3dmhvh2lz-go-containerregistry-0.19.1-crane ..."
		drvs := findDerivs(s)
		require.Equal(t, []string{"/nix/store/zhpnpc6ljvdp10ky2wfqwph3dmhvh2lz-go-containerregistry-0.19.1-crane"}, drvs)
	})
}
