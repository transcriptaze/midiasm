package metaevent

import (
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type CuePoint struct {
	Tag      string
	Status   types.Status
	Type     types.MetaEventType
	CuePoint string
}

func NewCuePoint(r io.ByteReader) (*CuePoint, error) {
	cuepoint, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &CuePoint{
		Tag:      "CuePoint",
		Status:   0xff,
		Type:     0x07,
		CuePoint: string(cuepoint),
	}, nil
}
