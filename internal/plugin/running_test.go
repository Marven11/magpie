package plugin

import (
	"testing"
	"time"
)

// Running is a boolean the Plugins page asks for on every look, and a host's
// startup holds hostMu for as long as Bun takes to download and the plugins
// to load (two minutes at most) — so a reader that takes that lock waits for
// a startup, and its own request being cancelled doesn't free it.
func TestRunningDoesNotWaitForAStartup(t *testing.T) {
	hostMu.Lock()
	locked := true
	defer func() {
		if locked {
			hostMu.Unlock()
		}
	}()
	done := make(chan bool, 1)
	go func() { done <- Running() }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Running() waited for a startup's lock: the Plugins page waits with it")
	}
	hostMu.Unlock()
	locked = false
}
