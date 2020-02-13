package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
)

type SequencerSpecificEvent struct {
	Tag          string
	Status       types.Status
	Type         types.MetaEventType
	Manufacturer types.Manufacturer
	Data         types.Hex
}

func NewSequencerSpecificEvent(ctx *context.Context, r events.EventReader, status types.Status, eventType types.MetaEventType) (*SequencerSpecificEvent, error) {
	if eventType != 0x7f {
		return nil, fmt.Errorf("Invalid SequencerSpecificEvent event type (%02x): expected '7F'", eventType)
	}

	bytes, err := r.ReadVLQ()
	if err != nil {
		return nil, err
	}

	id := bytes[0:1]
	data := bytes[1:]
	if bytes[0] == 0x00 {
		id = bytes[0:3]
		data = bytes[3:]
	}

	return &SequencerSpecificEvent{
		Tag:          "SequencerSpecificEvent",
		Status:       status,
		Type:         eventType,
		Manufacturer: ctx.LookupManufacturer(id),
		Data:         data,
	}, nil
}
