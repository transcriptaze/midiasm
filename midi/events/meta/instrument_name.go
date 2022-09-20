package metaevent

import (
	"github.com/transcriptaze/midiasm/midi/types"
)

type InstrumentName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewInstrumentName(bytes []byte) (*InstrumentName, error) {
	return &InstrumentName{
		Tag:    "InstrumentName",
		Status: 0xff,
		Type:   0x04,
		Name:   string(bytes),
	}, nil
}
