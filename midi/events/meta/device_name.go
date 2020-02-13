package metaevent

import (
	"github.com/twystd/midiasm/midi/types"
)

type DeviceName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewDeviceName(bytes []byte) (*DeviceName, error) {
	return &DeviceName{
		Tag:    "DeviceName",
		Status: 0xff,
		Type:   0x09,
		Name:   string(bytes),
	}, nil
}
