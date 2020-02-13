package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
)

type NoteOff struct {
	Tag      string
	Status   types.Status
	Channel  types.Channel
	Note     Note
	Velocity byte
}

func NewNoteOff(ctx *context.Context, r events.EventReader, status types.Status) (*NoteOff, error) {
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
		Note: Note{
			Value: note,
			Name:  ctx.GetNoteOff(channel, note),
			Alias: ctx.FormatNote(note),
		},
		Velocity: velocity,
	}, nil
}
