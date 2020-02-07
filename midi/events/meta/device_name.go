package metaevent

import (
	"fmt"
	"io"
)

type DeviceName struct {
	Tag string
	MetaEvent
	Name string
}

func NewDeviceName(event *MetaEvent, r io.ByteReader) (*DeviceName, error) {
	if event.Type != 0x09 {
		return nil, fmt.Errorf("Invalid DeviceName event type (%02x): expected '09'", event.Type)
	}

	name, err := read(r)
	if err != nil {
		return nil, err
	}

	return &DeviceName{
		Tag:       "DeviceName",
		MetaEvent: *event,
		Name:      string(name),
	}, nil
}
