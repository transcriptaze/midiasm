package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Copyright struct {
	Tag       string
	Status    types.Status
	Type      types.MetaEventType
	Copyright string
}

func NewCopyright(r io.ByteReader, status types.Status, eventType types.MetaEventType) (*Copyright, error) {
	if eventType != 0x02 {
		return nil, fmt.Errorf("Invalid Copyright event type (%02x): expected '02'", eventType)
	}

	data, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &Copyright{
		Tag:       "Copyright",
		Status:    status,
		Type:      eventType,
		Copyright: string(data),
	}, nil
}
