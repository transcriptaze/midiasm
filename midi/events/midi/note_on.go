package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type NoteOn struct {
	Tag string
	*events.Event
	Channel  types.Channel
	Note     Note
	Velocity byte
}

func NewNoteOn(ctx *context.Context, event *events.Event, r io.ByteReader) (*NoteOn, error) {
	if event.Status&0xF0 != 0x90 {
		return nil, fmt.Errorf("Invalid NoteOn status (%02x): expected '90'", event.Status&0x80)
	}

	channel := types.Channel((event.Status) & 0x0F)

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
		Event:   event,
		Channel: channel,
		Note: Note{
			Value: note,
			Name:  ctx.FormatNote(note),
			Alias: ctx.FormatNote(note),
		},
		Velocity: velocity,
	}, nil
}
