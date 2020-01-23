package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type SequencerSpecificEvent struct {
	MetaEvent
	Data types.Hex
}

func NewSequencerSpecificEvent(event *MetaEvent, r io.ByteReader) (*SequencerSpecificEvent, error) {
	if event.Type != 0x7f {
		return nil, fmt.Errorf("Invalid SequencerSpecificEvent event type (%02x): expected '7F'", event.Type)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	}

	return &SequencerSpecificEvent{
		MetaEvent: *event,
		Data:      data,
	}, nil
}
