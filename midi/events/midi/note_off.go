package midievent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	"github.com/transcriptaze/midiasm/midi/types"
)

type NoteOff struct {
	event
	Note     Note
	Velocity byte
}

func NewNoteOff(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status types.Status) (*NoteOff, error) {
	if status&0xf0 != 0x80 {
		return nil, fmt.Errorf("Invalid NoteOff status (%v): expected '8x'", status)
	}

	channel := types.Channel(status & 0x0f)

	note, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	velocity, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &NoteOff{
		event: event{
			tick:  tick,
			delta: delta,
			bytes: r.Bytes(),

			Tag:     "NoteOff",
			Status:  status,
			Channel: channel,
		},
		Note: Note{
			Value: note,
			Name:  GetNoteOff(ctx, channel, note),
			Alias: FormatNote(ctx, note),
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

	n.Note = Note{
		Value: note,
		Name:  ctx.GetNoteOff(n.Channel, note),
		Alias: FormatNote(ctx, note),
	}
}

func GetNoteOff(ctx *context.Context, ch types.Channel, n byte) string {
	if ctx != nil {
		return ctx.GetNoteOff(ch, n)
	}

	return FormatNote(ctx, n)
}
