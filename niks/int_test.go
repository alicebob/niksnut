package niks

import (
	"log/slog"
	"os"
	"strings"
	"testing"
)

func inInt() bool {
	return os.Getenv("INT") != ""
}

// Either makes a new tmp dir which can be used in tests, or skips the test if INT isn't set
// In verbose test (go test -v) the tmp dir isn't cleaned.
func SetupInt(t *testing.T) string {
	if !inInt() {
		t.Skip("integration test")
	}
	SetupLog(t)

	root, err := os.MkdirTemp("", "niksint")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("work dir: %s", root)
	t.Cleanup(func() {
		if testing.Verbose() {
			return
		}
		os.RemoveAll(root)
	})
	return root
}

// redirect slog output to testing.Log, which you only see on either error or "test -v"
func SetupLog(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler((*TestLog)(t), nil)))
}

type TestLog testing.T

func (t *TestLog) Write(l []byte) (int, error) {
	(*testing.T)(t).Log(strings.TrimSpace(string(l)))
	return 0, nil
}
