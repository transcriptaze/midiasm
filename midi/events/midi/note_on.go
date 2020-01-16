package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
)

type NoteOn struct {
	MidiEvent
	Note     Note
	Velocity byte
}

func NewNoteOn(ctx *context.Context, event *MidiEvent, r io.ByteReader) (*NoteOn, error) {
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
		Note: Note{
			Value: note,
			Name:  ctx.FormatNote(note),
		},
		Velocity: velocity,
	}, nil
}

func (e *NoteOn) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s channel:%-2v note:%-s, velocity:%d", e.MidiEvent, "NoteOn", e.Channel, e.Note.Name, e.Velocity)
}