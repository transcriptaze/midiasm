package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type MIDIChannelPrefix struct {
	Tag string
	*events.Event
	Type    types.MetaEventType
	Channel int8
}

func NewMIDIChannelPrefix(event *events.Event, eventType types.MetaEventType, r io.ByteReader) (*MIDIChannelPrefix, error) {
	if eventType != 0x20 {
		return nil, fmt.Errorf("Invalid MIDIChannelPrefix event type (%02x): expected '20'", eventType)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	} else if len(data) != 1 {
		return nil, fmt.Errorf("Invalid MIDIChannelPrefix length (%d): expected '1'", len(data))
	}

	channel := int8(data[0])
	if channel < 0 || channel > 15 {
		return nil, fmt.Errorf("Invalid MIDIChannelPrefix channel (%d): expected a value in the interval [0..15]", channel)
	}

	return &MIDIChannelPrefix{
		Tag:     "MIDIChannelPrefix",
		Event:   event,
		Type:    eventType,
		Channel: channel,
	}, nil
}
