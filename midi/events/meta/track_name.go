package metaevent

import (
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type TrackName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewTrackName(r io.ByteReader) (*TrackName, error) {
	name, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &TrackName{
		Tag:    "TrackName",
		Status: 0xff,
		Type:   0x03,
		Name:   string(name),
	}, nil
}
