package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
)

type MIDIPort struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Port   uint8
}

func NewMIDIPort(bytes []byte) (*MIDIPort, error) {
	if len(bytes) != 1 {
		return nil, fmt.Errorf("Invalid MIDIPort length (%d): expected '1'", len(bytes))
	}

	port := bytes[0]
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
