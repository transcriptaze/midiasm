package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/event"
	"io"
)

type NoteOff struct {
	MidiEvent
	Note     byte
	Velocity byte
}

func NewNoteOff(event *MidiEvent, r io.ByteReader) (*NoteOff, error) {
	if event.Status&0xF0 != 0x80 {
		return nil, fmt.Errorf("Invalid NoteOff status (%02x): expected '80'", event.Status&0xF0)
	}

	note, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	velocity, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &NoteOff{
		MidiEvent: *event,
		Note:      note,
		Velocity:  velocity,
	}, nil
}

func (e *NoteOff) Render(ctx *event.Context, w io.Writer) {
	note := ctx.Scale[e.Note%12]
	octave := -2 + int(e.Note)/12
	fmt.Fprintf(w, "%s %-16s channel:%d note:%-4s velocity:%d", e.MidiEvent, "NoteOff", e.Channel, fmt.Sprintf("%s%d", note, octave), e.Velocity)
}
