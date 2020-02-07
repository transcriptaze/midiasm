package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
)

type NoteOn struct {
	Tag string
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

	ctx.PutNoteOn(event.Channel, note)

	return &NoteOn{
		Tag:       "NoteOn",
		MidiEvent: *event,
		Note: Note{
			Value: note,
			Name:  ctx.FormatNote(note),
			Alias: ctx.FormatNote(note),
		},
		Velocity: velocity,
	}, nil
}
