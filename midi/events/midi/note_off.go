package midievent

import (
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/types"
)

type NoteOff struct {
	Tag      string
	Status   types.Status
	Channel  types.Channel
	Note     types.Note
	Velocity byte
}

func NewNoteOff(ctx *context.Context, r io.ByteReader, status types.Status) (*NoteOff, error) {
	if status&0xF0 != 0x80 {
		return nil, fmt.Errorf("Invalid NoteOff status (%v): expected '8x'", status)
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

	return &NoteOff{
		Tag:     "NoteOff",
		Status:  status,
		Channel: channel,
		Note: types.Note{
			Value: note,
			Name:  ctx.GetNoteOff(channel, note),
			Alias: ctx.FormatNote(note),
		},
		Velocity: velocity,
	}, nil
}

func (n *NoteOff) Transpose(ctx *context.Context, steps int) {
	v := int(n.Note.Value) + steps
	note := n.Note.Value

	switch {
	case v < 0:
		note = 0

	case v > 127:
		note = 127

	default:
		note = byte(v)
	}

	n.Note = types.Note{
		Value: note,
		Name:  ctx.GetNoteOff(n.Channel, note),
		Alias: ctx.FormatNote(note),
	}
}
