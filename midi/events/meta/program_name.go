package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
)

type ProgramName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewProgramName(r events.EventReader, status types.Status, eventType types.MetaEventType) (*ProgramName, error) {
	if eventType != 0x08 {
		return nil, fmt.Errorf("Invalid ProgramName event type (%02x): expected '08'", eventType)
	}

	name, err := r.ReadVLF()
	if err != nil {
		return nil, err
	}

	return &ProgramName{
		Tag:    "ProgramName",
		Status: status,
		Type:   eventType,
		Name:   string(name),
	}, nil
}
