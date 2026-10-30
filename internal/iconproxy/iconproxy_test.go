package iconproxy

import (
	"net/url"
	"testing"
)

func TestURLCarriesThePictureWhole(t *testing.T) {
	t.Setenv("MAGPIE_ICON_PROXY", "")
	pic := "https://github.com/magpie-community.png?size=96&x=a b#c"
	got := URL(pic)
	u, err := url.Parse(got)
	if err != nil {
		t.Fatal(err)
	}
	if u.Scheme+"://"+u.Host+u.Path != Base {
		t.Fatalf("URL(%q) = %q, not on %s", pic, got, Base)
	}
	if q := u.Query(); len(q) != 1 || q.Get("url") != pic {
		t.Fatalf("URL(%q) = %q: query %v", pic, got, q)
	}
}

func TestURLLeavesUsemagpieAlone(t *testing.T) {
	t.Setenv("MAGPIE_ICON_PROXY", "")
	for u, want := range map[string]string{
		"https://usemagpie.ai/partners/x.png": "https://usemagpie.ai/partners/x.png",
		// usemagpie.ai's own, spelled another way: the worker's own check
		// reads the hostname in lower case, so it refuses these too, and
		// asking it for one only sent it round to itself (a 400 to the app)
		"https://www.usemagpie.ai/x.png":          "https://www.usemagpie.ai/x.png",
		"https://cdn.usemagpie.ai/i/x.svg":        "https://cdn.usemagpie.ai/i/x.svg",
		"https://USEMAGPIE.AI/x.png":              "https://USEMAGPIE.AI/x.png",
		"https://usemagpie.ai:443/x.png":          "https://usemagpie.ai:443/x.png",
		"https://usemagpie.ai.evil.example/x.png": Base + "?url=" + url.QueryEscape("https://usemagpie.ai.evil.example/x.png"),
		"https://evil.example/?usemagpie.ai":      Base + "?url=" + url.QueryEscape("https://evil.example/?usemagpie.ai"),
		"https://evil.example/usemagpie.ai/x.png": Base + "?url=" + url.QueryEscape("https://evil.example/usemagpie.ai/x.png"),
	} {
		if got := URL(u); got != want {
			t.Errorf("URL(%q) = %q, want %q", u, got, want)
		}
	}
}
