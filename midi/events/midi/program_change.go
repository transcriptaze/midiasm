package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type ProgramChange struct {
	Tag string
	*events.Event
	Channel types.Channel
	Program byte
}

func NewProgramChange(event *events.Event, r io.ByteReader) (*ProgramChange, error) {
	if event.Status&0xF0 != 0xc0 {
		return nil, fmt.Errorf("Invalid ProgramChange status (%02x): expected 'C0'", event.Status&0x80)
	}

	channel := types.Channel((event.Status) & 0x0F)

	program, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &ProgramChange{
		Tag:     "ProgramChange",
		Event:   event,
		Channel: channel,
		Program: program,
	}, nil
}
