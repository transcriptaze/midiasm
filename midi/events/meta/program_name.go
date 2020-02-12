package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type ProgramName struct {
	Tag string
	MetaEvent
	Type types.MetaEventType
	Name string
}

func NewProgramName(event *MetaEvent, eventType types.MetaEventType, r io.ByteReader) (*ProgramName, error) {
	if eventType != 0x08 {
		return nil, fmt.Errorf("Invalid ProgramName event type (%02x): expected '08'", eventType)
	}

	name, err := read(r)
	if err != nil {
		return nil, err
	}

	return &ProgramName{
		Tag:       "ProgramName",
		MetaEvent: *event,
		Type:      eventType,
		Name:      string(name),
	}, nil
}
