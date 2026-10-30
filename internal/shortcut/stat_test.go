package shortcut

import (
	"os"
	"path/filepath"
	"testing"
)

// A Start menu that can't be looked at isn't one the user deleted anything
// from: a folder redirected to a network place, an ACL refusing this user.
// Remembering "made, and the user deleted it" there leaves a magpie that
// never makes the shortcut again, with no way back — and an upgrade from a
// magpie of before the state file is where it went wrong: that one
// registered itself at every start, so it read as the user having deleted it.
func TestEnsureShortcutItCannotLookAt(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "Start Menu")
	if err := os.WriteFile(file, []byte("a file where the folder would be"), 0o644); err != nil {
		t.Fatal(err)
	}
	// a path inside a file: Stat says ENOTDIR, which is not "there is none"
	lnk := filepath.Join(file, "Programs", Name+".lnk")
	state := filepath.Join(dir, "magpie", StateName)
	writes := 0
	oldMenu, oldWrite, oldReg, oldRegd := startMenu, writeLink, register, registered
	t.Cleanup(func() { startMenu, writeLink, register, registered = oldMenu, oldWrite, oldReg, oldRegd })
	startMenu = func() (string, error) { return lnk, nil }
	writeLink = func(lnk, exe string) (bool, error) { writes++; return true, nil }
	register = func(string) {}
	registered = func() bool { return true } // an older magpie registered itself
	for i := 0; i < 2; i++ {
		if err := ensure(`C:\Tools\magpie.exe`, "", state); err != nil {
			t.Fatal(err)
		}
	}
	if exists(state) {
		t.Error("a shortcut that can't be looked at was recorded as deleted by the user")
	}
	if writes != 0 {
		t.Error("a shortcut was written where nothing can be looked at")
	}
}
