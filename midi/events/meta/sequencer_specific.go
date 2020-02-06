package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type SequencerSpecificEvent struct {
	MetaEvent
	Manufacturer types.Manufacturer
	Data         types.Hex
}

func NewSequencerSpecificEvent(ctx *context.Context, event *MetaEvent, r io.ByteReader) (*SequencerSpecificEvent, error) {
	if event.Type != 0x7f {
		return nil, fmt.Errorf("Invalid SequencerSpecificEvent event type (%02x): expected '7F'", event.Type)
	}

	bytes, err := read(r)
	if err != nil {
		return nil, err
	}

	id := bytes[0:1]
	data := bytes[1:]
	if bytes[0] == 0x00 {
		id = bytes[0:3]
		data = bytes[3:]
	}

	return &SequencerSpecificEvent{
		MetaEvent:    *event,
		Manufacturer: ctx.LookupManufacturer(id),
		Data:         data,
	}, nil
}
