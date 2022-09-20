package metaevent

import (
	"github.com/transcriptaze/midiasm/midi/types"
)

type TrackName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewTrackName(bytes []byte) (*TrackName, error) {
	return &TrackName{
		Tag:    "TrackName",
		Status: 0xff,
		Type:   0x03,
		Name:   string(bytes),
	}, nil
}
