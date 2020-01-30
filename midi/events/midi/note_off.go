package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
)

type NoteOff struct {
	MidiEvent
	Note     Note
	Velocity byte
}

func NewNoteOff(ctx *context.Context, event *MidiEvent, r io.ByteReader) (*NoteOff, error) {
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
		Note: Note{
			Value: note,
			Name:  ctx.GetNoteOff(event.Channel, note),
			Alias: ctx.FormatNote(note),
		},
		Velocity: velocity,
	}, nil
}
