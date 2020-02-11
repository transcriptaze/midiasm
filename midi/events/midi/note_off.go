package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type NoteOff struct {
	Tag string
	*events.Event
	Channel  types.Channel
	Note     Note
	Velocity byte
}

func NewNoteOff(ctx *context.Context, event *events.Event, r io.ByteReader) (*NoteOff, error) {
	if event.Status&0xF0 != 0x80 {
		return nil, fmt.Errorf("Invalid NoteOff status (%02x): expected '80'", event.Status&0xF0)
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

	return &NoteOff{
		Tag:     "NoteOff",
		Event:   event,
		Channel: channel,
		Note: Note{
			Value: note,
			Name:  ctx.GetNoteOff(channel, note),
			Alias: ctx.FormatNote(note),
		},
		Velocity: velocity,
	}, nil
}
