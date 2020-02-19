package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type ProgramChange struct {
	Tag     string
	Status  types.Status
	Channel types.Channel
	Bank    uint16
	Program byte
}

func NewProgramChange(ctx *context.Context, r io.ByteReader, status types.Status) (*ProgramChange, error) {
	if status&0xF0 != 0xc0 {
		return nil, fmt.Errorf("Invalid ProgramChange status (%v): expected 'Cx'", status)
	}

	channel := uint8(status & 0x0f)

	program, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &ProgramChange{
		Tag:     "ProgramChange",
		Status:  status,
		Channel: types.Channel(channel),
		Bank:    ctx.ProgramBank[channel],
		Program: program,
	}, nil
}
