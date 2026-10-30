package cli

import (
	"os"
	"reflect"
	"testing"

	"github.com/yetone/magpie/internal/provider"
)

// idsCli is a provider to type model ids at.
func idsCli(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", os.Getenv("HOME"))
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	t.Setenv("DSH_HOME", "")
	if err := provider.Save(provider.Provider{ID: "relay", Name: "Relay", Key: "k1", Chat: "http://127.0.0.1:1/v1"}); err != nil {
		t.Fatal(err)
	}
}

// magpie models and magpie provider <id> print a model as <provider>/<model>,
// which is how an agent is told it too — so an id copied from either and typed
// back carries the provider's own prefix. Kept as given, the provider's models
// became relay/foo: a model the vendor doesn't have, under the provider's id
// twice, and the models it does serve were no longer exposed.
func TestProviderModelsTakesThePrefixItPrints(t *testing.T) {
	idsCli(t)
	if err := providerCmd([]string{"provider", "models", "relay", "relay/foo", "bar"}); err != nil {
		t.Fatal(err)
	}
	p, err := provider.Find("relay")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"foo", "bar"}; !reflect.DeepEqual(p.Models, want) {
		t.Fatalf("models %v, want %v", p.Models, want)
	}
	// adding one keeps the rest, as +id always did
	if err := providerCmd([]string{"provider", "models", "relay", "+relay/mini"}); err != nil {
		t.Fatal(err)
	}
	p, _ = provider.Find("relay")
	if want := []string{"foo", "bar", "mini"}; !reflect.DeepEqual(p.Models, want) {
		t.Fatalf("after +relay/mini: %v, want %v", p.Models, want)
	}
	// and models= takes it too
	if err := providerCmd([]string{"provider", "set", "relay", "models=relay/cheap,mini"}); err != nil {
		t.Fatal(err)
	}
	p, _ = provider.Find("relay")
	if want := []string{"cheap", "mini"}; !reflect.DeepEqual(p.Models, want) {
		t.Fatalf("models=: %v, want %v", p.Models, want)
	}
}
