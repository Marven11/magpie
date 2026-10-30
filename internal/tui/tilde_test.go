package tui

import (
	"path/filepath"
	"testing"
)

// One shortening to "~" for the paths a page shows: the sessions page had a
// second of its own that took /Users/cube-old for home /Users/cube and showed
// "~-old/x" for a folder that isn't under home at all, while the agents page
// showed the same path as it is.
func TestTildeIsOnlyForPathsUnderHome(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	if got, want := tilde(filepath.Join(homeDir, "work")), "~"+string(filepath.Separator)+"work"; got != want {
		t.Errorf("tilde(%q) = %q, want %q", filepath.Join(homeDir, "work"), got, want)
	}
	other := homeDir + "-old" + string(filepath.Separator) + "work"
	if got := tilde(other); got != other {
		t.Errorf("tilde(%q) = %q, want it unchanged: it isn't under home", other, got)
	}
}
