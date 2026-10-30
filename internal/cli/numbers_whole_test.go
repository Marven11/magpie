package cli

import (
	"os"
	"testing"
)

// A number typed with letters after it is not a number: Sscan took the part it
// read and dropped the rest, so "1abc" moved a search API and "3minutes" set a
// three-minute sync, and both said they did it.
func TestNumbersAreNumbersWhole(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", os.Getenv("HOME"))
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	for _, a := range [][]string{{"add", "tavily", "tvly-AAA"}, {"add", "brave", "brv-BBB"}} {
		if err := searchCmd(a); err != nil {
			t.Fatal(err)
		}
	}
	// the sync has to be on, or SetAuto answers "sync is off" whatever it is
	// given and the number never reaches it
	dav := newDAVServer(t, "pw")
	piped(t, "pw", "correct horse")
	if err := webdavCmd([]string{"on", dav.url, "user=me"}); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name string
		args []string
		run  func([]string) error
	}{
		{"search move", []string{"move", "brave", "1abc"}, searchCmd},
		{"search move, a word", []string{"move", "brave", "one"}, searchCmd},
		{"webdav auto", []string{"auto", "15x"}, webdavCmd},
		{"webdav auto unit", []string{"auto", "3minutes"}, webdavCmd},
	} {
		if err := c.run(c.args); err == nil {
			t.Errorf("%s %v was taken as a number", c.name, c.args)
		}
	}
	// and the numbers themselves still work
	if err := searchCmd([]string{"move", "brave", "1"}); err != nil {
		t.Fatal(err)
	}
	if err := webdavCmd([]string{"auto", "15"}); err != nil {
		t.Fatal(err)
	}
}
