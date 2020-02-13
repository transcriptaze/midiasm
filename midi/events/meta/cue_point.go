package metaevent

import (
	"github.com/twystd/midiasm/midi/types"
)

type CuePoint struct {
	Tag      string
	Status   types.Status
	Type     types.MetaEventType
	CuePoint string
}

func NewCuePoint(bytes []byte) (*CuePoint, error) {
	return &CuePoint{
		Tag:      "CuePoint",
		Status:   0xff,
		Type:     0x07,
		CuePoint: string(bytes),
	}, nil
}
