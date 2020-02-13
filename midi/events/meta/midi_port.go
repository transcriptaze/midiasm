package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
)

type MIDIPort struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Port   uint8
}

func NewMIDIPort(r events.EventReader, status types.Status, eventType types.MetaEventType) (*MIDIPort, error) {
	if eventType != 0x21 {
		return nil, fmt.Errorf("Invalid MIDIPort event type (%02x): expected '21'", eventType)
	}

	data, err := r.ReadVLF()
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
		Status: status,
		Type:   eventType,
		Port:   port & 0x7f,
	}, nil
}
