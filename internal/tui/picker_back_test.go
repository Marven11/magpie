package tui

import (
	"testing"

	"github.com/yetone/magpie/internal/provider"
)

// The picker's esc goes back to the mode it was opened from, and a routing
// group's pickers set that to the group page. providers' two pickers didn't
// set m.back, and goTo doesn't reset it: after a group's model picker had
// left modeGroup there, esc out of providers' add list or a provider's model
// list dropped into the group page — its keys then moved members, added
// models and changed routing, and saveGroup wrote providers.json, while the
// user was on the providers page.
func TestProvidersPickersLeaveBackAtTheirList(t *testing.T) {
	home(t)
	m := model{w: 120, h: 40}
	m.back = modeGroup // as a routing group's picker left it
	m.openPresets()
	if m.back != modeList {
		t.Fatalf("the add list's esc goes to mode %d, want the providers list", m.back)
	}
	if err := provider.Save(provider.Provider{ID: "acme", Name: "Acme", Chat: "https://acme.example.com/v1", Key: "k"}); err != nil {
		t.Fatal(err)
	}
	m = model{w: 120, h: 40}
	m.reload()
	m.back = modeGroup
	m.provs, m.prow = provider.All(), 0
	if len(m.provs) == 0 {
		t.Fatal("no providers loaded")
	}
	m.openProviderModels(m.provs[0].ID)
	if m.back != modeList {
		t.Fatalf("the model list's esc goes to mode %d, want the providers list", m.back)
	}
	// and esc through the key path lands on the list, not in a group
	next, _ := m.Update(keyMsg("esc"))
	m = next.(model)
	if m.mode != modeList {
		t.Fatalf("esc left mode %d", m.mode)
	}
}
