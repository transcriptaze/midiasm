package metaevent

import (
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type ProgramName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewProgramName(r io.ByteReader) (*ProgramName, error) {
	name, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &ProgramName{
		Tag:    "ProgramName",
		Status: 0xff,
		Type:   0x08,
		Name:   string(name),
	}, nil
}
