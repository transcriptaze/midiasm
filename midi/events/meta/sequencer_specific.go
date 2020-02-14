package metaevent

import (
	"github.com/twystd/midiasm/midi/types"
)

type SequencerSpecificEvent struct {
	Tag          string
	Status       types.Status
	Type         types.MetaEventType
	Manufacturer types.Manufacturer
	Data         types.Hex
}

func NewSequencerSpecificEvent(bytes []byte) (*SequencerSpecificEvent, error) {
	id := bytes[0:1]
	data := bytes[1:]
	if bytes[0] == 0x00 {
		id = bytes[0:3]
		data = bytes[3:]
	}

	return &SequencerSpecificEvent{
		Tag:          "SequencerSpecificEvent",
		Status:       0xff,
		Type:         0x7f,
		Manufacturer: types.LookupManufacturer(id),
		Data:         data,
	}, nil
}
