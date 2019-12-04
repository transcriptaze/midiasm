package context

import (
	"strconv"
	"testing"
)

func TestFormatNote(t *testing.T) {
	sharps := []string{"C", "C\u266f", "D", "D\u266f", "E", "F", "F\u266f", "G", "G\u266f", "A", "A\u266f", "B"}
	flats := []string{"C", "D\u266d", "D", "E\u266d", "E", "F", "G\u266d", "G", "A\u266d", "A", "B\u266d", "B"}

	ctx := Context{Scale: Sharps}
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

	ctx = Context{Scale: Flats}
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
