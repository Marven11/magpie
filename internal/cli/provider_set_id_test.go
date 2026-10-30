package cli

import (
	"strings"
	"testing"

	"github.com/yetone/magpie/internal/provider"
)

// A k=v pair the provider's own id can't take stops the command before
// anything is written: Save runs first and Rename checks the id after, so a
// typo in id= saved every other pair and then failed.
func TestProviderSetWritesNothingWhenTheIdIsBad(t *testing.T) {
	idsCli(t)
	for _, id := range []string{"", "My Relay", "with_underscore", "relay/other"} {
		err := providerCmd([]string{"provider", "set", "relay", "name=Renamed", "id=" + id})
		if err == nil || !strings.Contains(err.Error(), "lowercase letters, digits and dashes") {
			t.Fatalf("id=%q: %v", id, err)
		}
		p, _ := provider.Find("relay")
		if p.Name != "Relay" {
			t.Fatalf("id=%q: name %q was written before the id was refused", id, p.Name)
		}
	}
	// a good one still renames
	if err := providerCmd([]string{"provider", "set", "relay", "name=Renamed", "id=renamed"}); err != nil {
		t.Fatal(err)
	}
	if p, err := provider.Find("renamed"); err != nil || p.Name != "Renamed" {
		t.Fatalf("rename: %v %+v", err, p)
	}
}
