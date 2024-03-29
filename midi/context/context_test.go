package context

import (
	"strconv"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestFormatNote(t *testing.T) {
	SetMiddleC(lib.C3)

	sharps := []string{"C", "C\u266f", "D", "D\u266f", "E", "F", "F\u266f", "G", "G\u266f", "A", "A\u266f", "B"}
	flats := []string{"C", "D\u266d", "D", "E\u266d", "E", "F", "G\u266d", "G", "A\u266d", "A", "B\u266d", "B"}

	ctx := NewContext()
	note := byte(0)
	octave := -1
	for i := 0; i < 10; i++ {
		for j := 0; j < 12; j++ {
			expected := sharps[note%12] + strconv.Itoa(octave)
			if s := ctx.FormatNote(note); s != expected {
				t.Errorf("Invalid formatted note for %d (%02X): expected '%s', got '%s'", note, note, expected, s)
			}

			note += 1
		}
		octave += 1
	}

	ctx = NewContext().UseFlats()

	note = byte(0)
	octave = -1
	for i := 0; i < 10; i++ {
		for j := 0; j < 12; j++ {
			expected := flats[note%12] + strconv.Itoa(octave)
			if s := ctx.FormatNote(note); s != expected {
				t.Errorf("Invalid formatted note for %d (%02X): expected '%s', got '%s'", note, note, expected, s)
			}

			note += 1
		}
		octave += 1
	}
}

func TestFormatNoteWithC4(t *testing.T) {
	SetMiddleC(lib.C4)

	sharps := []string{"C", "C\u266f", "D", "D\u266f", "E", "F", "F\u266f", "G", "G\u266f", "A", "A\u266f", "B"}
	flats := []string{"C", "D\u266d", "D", "E\u266d", "E", "F", "G\u266d", "G", "A\u266d", "A", "B\u266d", "B"}

	ctx := NewContext()

	note := byte(0)
	octave := -2
	for i := 0; i < 10; i++ {
		for j := 0; j < 12; j++ {
			expected := sharps[note%12] + strconv.Itoa(octave)
			if s := ctx.FormatNote(note); s != expected {
				t.Errorf("Invalid formatted note for %d (%02X): expected '%s', got '%s'", note, note, expected, s)
			}

			note += 1
		}
		octave += 1
	}

	ctx = NewContext().UseFlats()

	note = byte(0)
	octave = -2
	for i := 0; i < 10; i++ {
		for j := 0; j < 12; j++ {
			expected := flats[note%12] + strconv.Itoa(octave)
			if s := ctx.FormatNote(note); s != expected {
				t.Errorf("Invalid formatted note for %d (%02X): expected '%s', got '%s'", note, note, expected, s)
			}

			note += 1
		}
		octave += 1
	}
}
