package davsync

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// A server that answers 403 for what isn't there (a Synology does for a
// shared folder not made) answers it for the usage folder of a computer that
// has no usage there yet, and the whole share round failed on it: the
// listing is the first call, and shareQuotas comes after it. The folder the
// usage one is in tells a 403 for what isn't there from one about access.
func TestDAVUsageListSeesAnUnmadeFolder(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u, p, ok := r.BasicAuth(); !ok || u != "me" || p != "pw" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch {
		case r.Method == "PROPFIND" && strings.TrimSuffix(r.URL.Path, "/") == "/magpie/magpie":
			w.WriteHeader(http.StatusMultiStatus)
		case strings.HasPrefix(r.URL.Path, "/magpie/magpie/usage"):
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusForbidden)
		}
	}))
	defer srv.Close()
	d, err := newDAV(Config{URL: srv.URL + "/magpie", User: "me", Password: "pw"})
	if err != nil {
		t.Fatal(err)
	}
	got, err := d.list(context.Background())
	if err != nil {
		t.Fatalf("listing the usage folder that isn't there: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("files: %v", got)
	}

	// a folder that isn't there either is still said: there is nothing to
	// sync into
	d, err = newDAV(Config{URL: srv.URL + "/missing", User: "me", Password: "pw"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := d.list(context.Background()); err == nil || !strings.Contains(err.Error(), "/missing") {
		t.Fatalf("listing a folder that isn't there: %v", err)
	}
}
