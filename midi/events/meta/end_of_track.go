package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type EndOfTrack struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
}

func NewEndOfTrack(r io.ByteReader) (*EndOfTrack, error) {
	data, err := events.VLF(r)
	if err != nil {
		return nil, err
	} else if len(data) != 0 {
		return nil, fmt.Errorf("Invalid EndOfTrack length (%d): expected '0'", len(data))
	}

	return &EndOfTrack{
		Tag:    "EndOfTrack",
		Status: 0xff,
		Type:   0x2f,
	}, nil
}
