package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/yetone/magpie/internal/provider"
)

// One value, and one value only: the words after it were dropped and the
// command said it did what was asked.
func TestModelTakesOneValue(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", os.Getenv("HOME"))
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	if err := provider.Save(provider.Provider{ID: "relay", Name: "Relay", Key: "k1", Chat: "http://127.0.0.1:1/v1"}); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		args []string
		run  func([]string) error
	}{
		{[]string{"context", "relay/foo", "1m", "bogus"}, modelCmd},
		{[]string{"output", "relay/foo", "4k", "bogus"}, modelCmd},
	} {
		err := c.run(c.args)
		if err == nil || !strings.Contains(err.Error(), "bogus") {
			t.Errorf("%v: %v", c.args, err)
		}
	}
}
