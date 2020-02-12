package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Copyright struct {
	Tag string
	MetaEvent
	Type      types.MetaEventType
	Copyright string
}

func NewCopyright(event *MetaEvent, eventType types.MetaEventType, r io.ByteReader) (*Copyright, error) {
	if eventType != 0x02 {
		return nil, fmt.Errorf("Invalid Copyright event type (%02x): expected '02'", eventType)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	}

	return &Copyright{
		Tag:       "Copyright",
		MetaEvent: *event,
		Type:      eventType,
		Copyright: string(data),
	}, nil
}
