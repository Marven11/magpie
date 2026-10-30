package catalog

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A cache file written over — a multi-megabyte write cut short, by another
// magpie reading it mid-write or by a process killed during one — was read as
// "no catalog" and latched: loaded was set before the read, so this process
// never read the file again, and Source/Stale/Missing went on calling it a
// fresh catalog. Provider answered nothing, every model priced at zero, for a
// day. The same file written again is read.
func TestCatalogReadsTheCacheAgainAfterAFailedRead(t *testing.T) {
	p := CachePath()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Remove(p)
		Reset()
	})
	Reset()
	// half a JSON document, as O_TRUNC-then-write leaves it
	if err := os.WriteFile(p, []byte(`{"openai":{"id":"openai","models":{"gpt-5":{"id":"gpt-5"`), 0o644); err != nil {
		t.Fatal(err)
	}
	if m := load(); len(m) != 0 {
		t.Fatalf("a truncated cache answered %d providers", len(m))
	}
	// the write finished
	full := `{"openai":{"id":"openai","models":{"gpt-5":{"id":"gpt-5","name":"GPT-5","limit":{"context":400000}}}}}`
	if err := os.WriteFile(p, []byte(full), 0o644); err != nil {
		t.Fatal(err)
	}
	if m := load(); len(m) == 0 {
		t.Fatal("the cache was written again and still read as none")
	}
}

// A Gemini list that still has a next page after the cap is not the list: read
// as one, the models past the cap look withdrawn to Refetch, which drops what
// the user picked of them. Paged (live.go) refuses the same shape.
func TestFetchGeminiTakesEndlessPagingAsNoAnswer(t *testing.T) {
	var asked int
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		asked++
		fmt.Fprintf(rw, `{"models":[{"name":"models/m-%d","supportedGenerationMethods":["generateContent"]}],"nextPageToken":"t%d"}`, asked, asked)
	}))
	defer srv.Close()
	if _, _, err := FetchGemini(context.Background(), srv.URL+"/v1beta", "k", nil); err == nil {
		t.Fatal("a list that never ends was read as the list")
	} else if !strings.Contains(err.Error(), "still paging after 20 pages") {
		t.Errorf("error %q, want it to say the paging didn't end", err)
	}
	if asked > 20 {
		t.Errorf("asked %d pages, want at most 20", asked)
	}
	// one that ends at the second page is the list
	asked = 0
	short := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		asked++
		if asked == 1 {
			rw.Write([]byte(`{"models":[{"name":"models/a","supportedGenerationMethods":["generateContent"]}],"nextPageToken":"t2"}`))
			return
		}
		rw.Write([]byte(`{"models":[{"name":"models/b","supportedGenerationMethods":["generateContent"]}]}`))
	}))
	defer short.Close()
	out, _, err := FetchGemini(context.Background(), short.URL+"/v1beta", "k", nil)
	if err != nil || len(out) != 2 {
		t.Fatalf("models %v, %v", out, err)
	}
}
