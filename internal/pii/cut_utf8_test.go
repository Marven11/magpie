package pii

import (
	"strings"
	"testing"
)

// A request's bytes are not checked for UTF-8 before they are cut: redact
// hands a JSON string's bytes over as they came when the string holds no
// backslash. A piece of such a text whose second half holds no line end and no
// space has no rune start to begin the next piece at, and the walk to one ran
// past the text — index out of range, in the goroutine masking.find starts, so
// the process went with it.
func TestPiecesCutWhatIsNotUTF8(t *testing.T) {
	for _, size := range []int{pieceSize, guardChunk} {
		for _, text := range []string{
			strings.Repeat("\x80", 2*size+3),
			"\n" + strings.Repeat("\x80", size+size/2+7),
			strings.Repeat("a", size-1024) + strings.Repeat("\x80", size),
		} {
			spans := cut(text, size)
			if len(spans) == 0 {
				t.Fatalf("size %d, length %d: no piece", size, len(text))
			}
			if spans[0].at != 0 || spans[len(spans)-1].end != len(text) {
				t.Fatalf("size %d: pieces %v don't cover the text of %d bytes", size, spans, len(text))
			}
			for i, s := range spans {
				if s.end <= s.at || s.end-s.at > size {
					t.Fatalf("size %d: piece %d is %v, of %d bytes at most", size, i, s, size)
				}
				if i > 0 && s.at > spans[i-1].end {
					t.Fatalf("size %d: piece %d starts at %d, after %d, so bytes went", size, i, s.at, spans[i-1].end)
				}
			}
		}
	}
}
