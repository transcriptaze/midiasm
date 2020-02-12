package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Text struct {
	Tag string
	MetaEvent
	Type types.MetaEventType
	Text string
}

func NewText(event *MetaEvent, eventType types.MetaEventType, r io.ByteReader) (*Text, error) {
	if eventType != 0x01 {
		return nil, fmt.Errorf("Invalid Text event type (%02x): expected '01'", eventType)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	}

	return &Text{
		Tag:       "Text",
		MetaEvent: *event,
		Type:      eventType,
		Text:      string(data),
	}, nil
}
