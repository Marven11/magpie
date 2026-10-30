//go:build !windows

package maskrec

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

// A day is parsed with no lock held: Hold, on a request's way, takes
// store.mu, and a day's file is tens of megabytes of JSONL — the privacy page
// asks for the summary and a page every 15s, which held every masked request
// for as long as the read took. The day here is a named pipe, so its read
// stays open until this writes to it: with store.mu held by the test, List
// gets to that open only if it doesn't need the lock to read the day.
func TestListReadsTheDayWithoutTheLock(t *testing.T) {
	now := time.Date(2026, 10, 20, 15, 0, 0, 0, time.Local)
	dir := sandbox(t, now)
	fifo := filepath.Join(dir, now.Format(dayForm)+".jsonl")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	store.mu.Lock()
	done := make(chan []Record, 1)
	go func() {
		rs, _ := List(0, 10, Filter{})
		done <- rs
	}()
	// a writer's open waits for a reader's: the one List takes of the day
	written := make(chan error, 1)
	go func() {
		f, err := os.OpenFile(fifo, os.O_WRONLY, 0o600)
		if err != nil {
			written <- err
			return
		}
		_, err = f.WriteString(`{"key":1,"t":"2026-10-20T15:00:00Z","what":"chat","agent":"codex","n":1}` + "\n")
		f.Close()
		written <- err
	}()
	select {
	case err := <-written:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(30 * time.Second):
		store.mu.Unlock()
		t.Fatal("List holds store.mu while it reads the day")
	}
	store.mu.Unlock()
	if rs := <-done; len(rs) != 1 || rs[0].Key != 1 {
		t.Fatalf("records: %+v", rs)
	}
}
