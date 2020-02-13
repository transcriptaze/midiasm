package metaevent

import (
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type DeviceName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewDeviceName(r io.ByteReader) (*DeviceName, error) {
	name, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &DeviceName{
		Tag:    "DeviceName",
		Status: 0xff,
		Type:   0x09,
		Name:   string(name),
	}, nil
}
