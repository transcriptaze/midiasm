package metaevent

import (
	"fmt"
	"io"
)

type Copyright struct {
	MetaEvent
	Copyright string
}

func NewCopyright(event *MetaEvent, r io.ByteReader) (*Copyright, error) {
	if event.Type != 0x02 {
		return nil, fmt.Errorf("Invalid Copyright event type (%02x): expected '02'", event.Type)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	}

	return &Copyright{
		MetaEvent: *event,
		Copyright: string(data),
	}, nil
}
