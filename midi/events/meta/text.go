package metaevent

import (
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Text struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Text   string
}

func NewText(r io.ByteReader) (*Text, error) {
	data, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &Text{
		Tag:    "Text",
		Status: 0xff,
		Type:   0x01,
		Text:   string(data),
	}, nil
}
