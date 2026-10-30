package plugin

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// A plugin-auth.json.lock that can't be taken or deleted — a host died
// holding it in a folder that refuses this user — used to spin at 200 tries a
// second with no way out, so SignOut (the GUI's sign-out, forgetting an
// account, the CLI's) never returned, and its caller's cancellation didn't
// end it either.
func TestLockAuthGivesUpWhenTheContextDoes(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if err := os.MkdirAll(filepath.Dir(AuthPath()), 0o755); err != nil {
		t.Fatal(err)
	}
	lock := AuthPath() + ".lock"
	if err := os.WriteFile(lock, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	began := time.Now()
	unlock, err := lockAuth(ctx)
	if err == nil {
		if unlock != nil {
			unlock()
		}
		t.Fatal("a lock another process holds was taken anyway")
	}
	if !errors.Is(err, context.DeadlineExceeded) || !strings.Contains(err.Error(), "can't be taken") {
		t.Fatalf("lockAuth = %v, want the lock named and the deadline in it", err)
	}
	if took := time.Since(began); took > 5*time.Second {
		t.Fatalf("it took %s to give up", took)
	}
}
