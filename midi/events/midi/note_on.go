package midievent

import (
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/types"
)

type NoteOn struct {
	Tag      string
	Status   types.Status
	Channel  types.Channel
	Note     types.Note
	Velocity byte
}

func NewNoteOn(ctx *context.Context, r io.ByteReader, status types.Status) (*NoteOn, error) {
	if status&0xF0 != 0x90 {
		return nil, fmt.Errorf("Invalid NoteOn status (%v): expected '9x'", status)
	}

	channel := types.Channel(status & 0x0F)

	note, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	velocity, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	ctx.PutNoteOn(channel, note)

	return &NoteOn{
		Tag:     "NoteOn",
		Status:  status,
		Channel: channel,
		Note: types.Note{
			Value: note,
			Name:  ctx.FormatNote(note),
			Alias: ctx.FormatNote(note),
		},
		Velocity: velocity,
	}, nil
}
