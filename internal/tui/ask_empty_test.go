package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// An empty answer isn't dropped on the floor: the step's own take on it —
// ParseRule's "a rule needs use=…", profile.Save's "profile name is empty" —
// is what says so in the status line, where a silent return read as a key that
// doesn't work.
func TestAnEmptyAnswerIsNotSwallowed(t *testing.T) {
	took := false
	m := model{w: 120, h: 40}
	m.mode, m.back = modeAsk, modeList
	m.ask = ask{input: newInput("the rule"), onEnter: func(v string) tea.Cmd {
		took = true
		return nil
	}}
	next, _ := m.Update(keyMsg("enter"))
	m = next.(model)
	if !took {
		t.Fatal("an empty answer was dropped without the step seeing it")
	}
	if m.mode != modeList {
		t.Fatalf("mode %d, want the prompt closed", m.mode)
	}
}
