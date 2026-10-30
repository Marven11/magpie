package redact

import (
	"os"
	"path/filepath"
	"testing"
)

// The key file is written the way magpie's other state files are — flushed to a
// temporary file and renamed over — so a reader after a crash sees the whole
// key or none, never 0 bytes: an empty file is read as "no key", and a new key
// renames every placeholder a vendor's cached conversation was primed with.
func TestKeyFileIsRenamedIn(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "redact.key")
	if err := os.WriteFile(p, nil, 0o600); err != nil { // a write that left nothing
		t.Fatal(err)
	}
	before, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	keyOf, keyPath = nil, ""
	mu.Unlock()
	SetKeyPath(p)
	k := key()
	after, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	if os.SameFile(before, after) {
		t.Error("the key was written into the file in place, so a crash leaves half of it")
	}
	b, err := os.ReadFile(p)
	if err != nil || len(b) != 32 || string(b) != string(k) {
		t.Fatalf("the file holds %d bytes, err %v; the key is %x", len(b), err, k)
	}
	if len(k) != 32 {
		t.Fatalf("key %x", k)
	}
	// and it is the key of the next run, which is what the file is for
	mu.Lock()
	keyOf = nil
	mu.Unlock()
	if again := key(); string(again) != string(k) {
		t.Error("the key changed across a restart: every placeholder does with it")
	}
}
