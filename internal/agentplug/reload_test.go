package agentplug

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yetone/magpie/internal/plugin"
	"github.com/yetone/magpie/internal/settings"
)

// use is plugins.json with these entries, as Settings › Plugins writes it.
func use(t *testing.T, entries ...plugin.Entry) {
	t.Helper()
	if err := os.MkdirAll(settings.Dir(), 0o755); err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(plugin.List{Plugins: entries})
	if err := os.WriteFile(filepath.Join(settings.Dir(), "plugins.json"), b, 0o644); err != nil {
		t.Fatal(err)
	}
	Reload()
	t.Cleanup(func() {
		os.Remove(filepath.Join(settings.Dir(), "plugins.json"))
		Reload()
	})
}

// An agent module bundled out of its entry and a file beside it (what a path
// plugin with a package.json is) was compiled once and never again: stamp()
// looked at the entry file alone, so a change to the file it imports went
// unseen and the old one kept being listed.
func TestListReloadsAnImportedModule(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"magpie":{"agent":"x.agent.js"}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	rules := filepath.Join(dir, "rules.js")
	if err := os.WriteFile(rules, []byte(`export const model = "one";`), 0o644); err != nil {
		t.Fatal(err)
	}
	entry := filepath.Join(dir, "x.agent.js")
	src := "import { model } from \"./rules.js\";\nexport const agent = { id: \"reload\", config: \"~/.reload.json\", model: \"model\" };\n" +
		"export function connect() { return { model }; }\n"
	if err := os.WriteFile(entry, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	use(t, plugin.Entry{Spec: dir})
	list, errs := List()
	if len(list) != 1 || len(errs) != 0 || list[0].Desc.ID != "reload" {
		t.Fatalf("agents %v, errs %v", list, errs)
	}
	// the module changes, the entry does not
	if err := os.WriteFile(rules, []byte(`export const model = "two";`), 0o644); err != nil {
		t.Fatal(err)
	}
	nextCheck.Store(0)
	list, errs = List()
	if len(errs) != 0 || len(list) != 1 {
		t.Fatalf("agents %v, errs %v", list, errs)
	}
	got, err := list[0].Connect(nil)
	if err != nil {
		t.Fatal(err)
	}
	if got["model"] != "two" {
		t.Fatalf("connect returned %v, want the changed module's model", got)
	}
}
