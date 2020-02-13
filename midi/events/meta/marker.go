package metaevent

import (
	"github.com/twystd/midiasm/midi/types"
)

type Marker struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Marker string
}

func NewMarker(bytes []byte) (*Marker, error) {
	return &Marker{
		Tag:    "Marker",
		Status: 0xff,
		Type:   0x06,
		Marker: string(bytes),
	}, nil
}
