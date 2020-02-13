package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type MIDIPort struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Port   uint8
}

func NewMIDIPort(r io.ByteReader) (*MIDIPort, error) {
	data, err := events.VLF(r)
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
		Tag:    "MIDIPort",
		Status: 0xff,
		Type:   0x21,
		Port:   port & 0x7f,
	}, nil
}
