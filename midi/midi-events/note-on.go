package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/event"
	"io"
)

type NoteOn struct {
	MidiEvent
	Note     byte
	Velocity byte
}

func NewNoteOn(event *MidiEvent, r io.ByteReader) (*NoteOn, error) {
	if event.Status&0xF0 != 0x90 {
		return nil, fmt.Errorf("Invalid NoteOn status (%02x): expected '90'", event.Status&0x80)
	}

	note, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	velocity, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &NoteOn{
		MidiEvent: *event,
		Note:      note,
		Velocity:  velocity,
	}, nil
}

func (e *NoteOn) Render(ctx *event.Context, w io.Writer) {
	note := ctx.Scale[e.Note%12]
	octave := -2 + int(e.Note)/12
	fmt.Fprintf(w, "%s %-16s channel:%d note:%-4s velocity:%d", e.MidiEvent, "NoteOn", e.Channel, fmt.Sprintf("%s%d", note, octave), e.Velocity)
}
