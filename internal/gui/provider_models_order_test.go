package gui

import (
	"os"
	"slices"
	"testing"

	"github.com/yetone/magpie/internal/catalog"
	"github.com/yetone/magpie/internal/provider"
)

// Picks the vendor's list doesn't know are shown first, the last one
// picked foremost, then the listed models in the list's order.
func TestUnlistedPicksShownFirst(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", os.Getenv("HOME"))
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	p := provider.Provider{ID: "relay", Name: "Relay", Chat: "https://relay.example/v1"}
	ms := func(ids ...string) (out []catalog.Model) {
		for _, id := range ids {
			out = append(out, catalog.Model{ID: id})
		}
		return
	}
	info := providerInfoWith(p, provider.Listing{Available: ms("a", "b"), Exposed: ms("a", "x", "y", "z")}, nil, nil)
	var ids []string
	for _, m := range info.Models {
		ids = append(ids, m.ID)
	}
	if want := []string{"z", "y", "x", "a", "b"}; !slices.Equal(ids, want) {
		t.Fatalf("models shown as %v, want %v", ids, want)
	}
	if info.Exposed != 4 {
		t.Fatalf("%d exposed, want 4", info.Exposed)
	}
}
