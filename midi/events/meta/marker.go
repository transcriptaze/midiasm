package metaevent

import (
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Marker struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Marker string
}

func NewMarker(r io.ByteReader) (*Marker, error) {
	marker, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &Marker{
		Tag:    "Marker",
		Status: 0xff,
		Type:   0x06,
		Marker: string(marker),
	}, nil
}
