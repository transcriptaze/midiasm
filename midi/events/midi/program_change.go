package midievent

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/types"
)

type ProgramChange struct {
	event
	Bank    uint16
	Program byte
}

func NewProgramChange(ctx *context.Context, tick uint64, delta uint32, r io.ByteReader, status types.Status) (*ProgramChange, error) {
	if status&0xF0 != 0xc0 {
		return nil, fmt.Errorf("Invalid ProgramChange status (%v): expected 'Cx'", status)
	}

	channel := uint8(status & 0x0f)

	program, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &ProgramChange{
		event: event{
			tick:  tick,
			delta: delta,
			bytes: []byte{0x00, byte(status), program},

			Tag:     "ProgramChange",
			Status:  status,
			Channel: types.Channel(channel),
		},
		Bank:    ctx.ProgramBank[channel],
		Program: program,
	}, nil
}

func (p ProgramChange) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0xc0 | p.Channel),
		byte(p.Program),
	}

	return
}

func (p *ProgramChange) UnmarshalText(bytes []byte) error {
	p.tick = 0
	p.delta = 0
	p.bytes = []byte{}
	p.Tag = "ProgramChange"

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)ProgramChange\s+channel:([0-9]+)\s+bank:([0-9]+),\s*program:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 5 {
		return fmt.Errorf("invalid ProgramChange event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if channel, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if bank, err := strconv.ParseUint(match[3], 10, 16); err != nil {
		return err
	} else if program, err := strconv.ParseUint(match[4], 10, 8); err != nil {
		return err
	} else if channel > 15 {
		return fmt.Errorf("invalid ProgramChange channel (%v)", channel)
	} else {
		p.delta = uint32(delta)
		p.bytes = []byte{0x00, byte(0xc0 | uint8(channel&0x0f)), byte(program)}
		p.Status = types.Status(0xc0 | uint8(channel&0x0f))
		p.Channel = types.Channel(channel)
		p.Bank = uint16(bank)
		p.Program = uint8(program)
	}

	return nil
}
