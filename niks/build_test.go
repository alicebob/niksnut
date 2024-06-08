package niks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindDerivs(t *testing.T) {
	s := "\t\t\t\techo that was it!.\n\t\t\t\techo pwd: $(pwd)\n\t\t\t\techo readlink: $(readlink -f ./result/)\n\t\t\t\techo result: $(ls ./result/)\n\t\t\t\techo ENV: $(printenv)\n\t\t\t\t$(/nix/store/0g1s8yd0biawp32fl3i7kdbi219jx6aq-openssh-9.7p1/bin/ssh --version)\n\t\t\t"
	drvs := findDerivs(s)
	require.Equal(t, []string{"/nix/store/0g1s8yd0biawp32fl3i7kdbi219jx6aq-openssh-9.7p1"}, drvs)
}
