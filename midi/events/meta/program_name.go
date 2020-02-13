package metaevent

import (
	"github.com/twystd/midiasm/midi/types"
)

type ProgramName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewProgramName(bytes []byte) (*ProgramName, error) {
	return &ProgramName{
		Tag:    "ProgramName",
		Status: 0xff,
		Type:   0x08,
		Name:   string(bytes),
	}, nil
}
