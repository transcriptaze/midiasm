package metaevent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/types"
)

type MIDIChannelPrefix struct {
	Tag     string
	Status  types.Status
	Type    types.MetaEventType
	Channel int8
}

func NewMIDIChannelPrefix(bytes []byte) (*MIDIChannelPrefix, error) {
	if len(bytes) != 1 {
		return nil, fmt.Errorf("Invalid MIDIChannelPrefix length (%d): expected '1'", len(bytes))
	}

	channel := int8(bytes[0])
	if channel < 0 || channel > 15 {
		return nil, fmt.Errorf("Invalid MIDIChannelPrefix channel (%d): expected a value in the interval [0..15]", channel)
	}

	return &MIDIChannelPrefix{
		Tag:     "MIDIChannelPrefix",
		Status:  0xff,
		Type:    0x20,
		Channel: channel,
	}, nil
}
