package tui

import (
	"strings"
	"testing"
)

// $VISUAL and $EDITOR are a command line, not a word: a program whose path
// holds a space comes quoted, and splitting on spaces alone named a program
// that isn't there (stat /Applications/Sublime: no such file).
func TestEditorArgs(t *testing.T) {
	for _, c := range []struct {
		in   string
		want []string
	}{
		{"/usr/bin/vi", []string{"/usr/bin/vi"}},
		{"code -w", []string{"code", "-w"}},
		{`"/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl" -w`,
			[]string{"/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl", "-w"}},
		{`'/opt/My Editor/ed' --wait`, []string{"/opt/My Editor/ed", "--wait"}},
		{`code --goto "a b.go"`, []string{"code", "--goto", "a b.go"}},
		{"  ", nil},
	} {
		if got := editorArgs(c.in); strings.Join(got, "|") != strings.Join(c.want, "|") {
			t.Errorf("editorArgs(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
