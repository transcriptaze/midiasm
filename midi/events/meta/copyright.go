package metaevent

import (
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

func NewCopyright(r io.ByteReader) (*Copyright, error) {
	data, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &Copyright{
		Tag:       "Copyright",
		Status:    0xff,
		Type:      0x02,
		Copyright: string(data),
	}, nil
}
