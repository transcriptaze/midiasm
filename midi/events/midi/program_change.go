package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type ProgramChange struct {
	Tag     string
	Status  types.Status
	Channel types.Channel
	Program byte
}

func NewProgramChange(r io.ByteReader, status types.Status) (*ProgramChange, error) {
	if status&0xF0 != 0xc0 {
		return nil, fmt.Errorf("Invalid ProgramChange status (%v): expected 'Cx'", status)
	}

	channel := types.Channel(status & 0x0F)

	program, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &ProgramChange{
		Tag:     "ProgramChange",
		Status:  status,
		Channel: channel,
		Program: program,
	}, nil
}
