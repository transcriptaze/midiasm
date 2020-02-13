package metaevent

import (
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type InstrumentName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewInstrumentName(r io.ByteReader) (*InstrumentName, error) {
	name, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &InstrumentName{
		Tag:    "InstrumentName",
		Status: 0xff,
		Type:   0x04,
		Name:   string(name),
	}, nil
}
