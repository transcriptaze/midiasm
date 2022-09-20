package metaevent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/types"
)

type EndOfTrack struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
}

func NewEndOfTrack(bytes []byte) (*EndOfTrack, error) {
	if len(bytes) != 0 {
		return nil, fmt.Errorf("Invalid EndOfTrack length (%d): expected '0'", len(bytes))
	}

	return &EndOfTrack{
		Tag:    "EndOfTrack",
		Status: 0xff,
		Type:   0x2f,
	}, nil
}
