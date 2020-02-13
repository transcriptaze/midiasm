package metaevent

import (
	"github.com/twystd/midiasm/midi/types"
)

type Copyright struct {
	Tag       string
	Status    types.Status
	Type      types.MetaEventType
	Copyright string
}

func NewCopyright(bytes []byte) (*Copyright, error) {
	return &Copyright{
		Tag:       "Copyright",
		Status:    0xff,
		Type:      0x02,
		Copyright: string(bytes),
	}, nil
}
