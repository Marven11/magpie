// Package jsmod compiles a plugin's JavaScript module for moejs, the
// JavaScript engine written in Go that magpie runs plugin code in its own
// process with: gateway middleware (internal/middleware) and agents
// (internal/agentplug).
package jsmod

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MoeForge/moejs"
)

// Compile compiles the module in file and the modules it imports, which
// are files under root: there is no node_modules here, so a package is
// bundled into the module. what names the kind of module in an error
// ("a middleware").
func Compile(root, file, what string) (*moejs.Module, error) {
	m, _, err := CompileFiles(root, file, what)
	return m, err
}

// CompileFiles is Compile and every file it read: file and the modules it
// imports, so a caller can tell when one of them changes. On an error the
// files read up to it are returned too.
func CompileFiles(root, file, what string) (*moejs.Module, []string, error) {
	src, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}
	files := []string{file}
	entry, err := moejs.Compile(file, string(src))
	if err != nil {
		return nil, files, err
	}
	mods := map[string]*moejs.Module{file: entry}
	mod, err := moejs.Link(entry, func(ref moejs.Referrer, spec string) (*moejs.Module, error) {
		if !strings.HasPrefix(spec, "./") && !strings.HasPrefix(spec, "../") {
			return nil, fmt.Errorf("%q: %s imports only files beside it; bundle packages into it", spec, what)
		}
		from := file
		if ref != nil {
			from = ref.Name()
		}
		p := filepath.Clean(filepath.Join(filepath.Dir(from), filepath.FromSlash(spec)))
		if rel, err := filepath.Rel(root, p); err != nil || strings.HasPrefix(rel, "..") {
			return nil, fmt.Errorf("%q is outside the plugin", spec)
		}
		if m, ok := mods[p]; ok {
			return m, nil
		}
		src, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		m, err := moejs.Compile(p, string(src))
		if err != nil {
			return nil, err
		}
		mods[p] = m
		files = append(files, p)
		return m, nil
	})
	return mod, files, err
}
