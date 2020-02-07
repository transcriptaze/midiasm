package metaevent

import (
	"fmt"
	"io"
)

type MIDIPort struct {
	Tag string
	MetaEvent
	Port uint8
}

func NewMIDIPort(event *MetaEvent, r io.ByteReader) (*MIDIPort, error) {
	if event.Type != 0x21 {
		return nil, fmt.Errorf("Invalid MIDIPort event type (%02x): expected '21'", event.Type)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	}

	if len(data) != 1 {
		return nil, fmt.Errorf("Invalid MIDIPort length (%d): expected '1'", len(data))
	}

	port := data[0]
	if port > 127 {
		return nil, fmt.Errorf("Invalid MIDIPort port (%d): expected a value in the interval [0..127]", port)
	}

	return &MIDIPort{
		Tag:       "MIDIPort",
		MetaEvent: *event,
		Port:      port & 0x7f,
	}, nil
}
