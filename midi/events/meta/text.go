package metaevent

import (
	"github.com/twystd/midiasm/midi/types"
)

type Text struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Text   string
}

func NewText(bytes []byte) (*Text, error) {
	return &Text{
		Tag:    "Text",
		Status: 0xff,
		Type:   0x01,
		Text:   string(bytes),
	}, nil
}
