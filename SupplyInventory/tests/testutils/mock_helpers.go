package testutils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// SetEnvAndRestore sets an env var and returns a function to restore previous value
func SetEnvAndRestore(key, value string) func() {
	old, had := os.LookupEnv(key)
	_ = os.Setenv(key, value)
	return func() {
		if had {
			_ = os.Setenv(key, old)
		} else {
			_ = os.Unsetenv(key)
		}
	}
}

// CreateTempDir creates a temporary directory for tests and returns its path and a cleanup func
func CreateTempDir(t *testing.T, prefix string) (string, func()) {
	t.Helper()
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return filepath.Clean(dir), func() { _ = os.RemoveAll(dir) }
}
