package midievent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	lib "github.com/transcriptaze/midiasm/midi/types"
)

type ProgramChange struct {
	event
	Bank    uint16
	Program byte
}

func MakeProgramChange(tick uint64, delta uint32, channel lib.Channel, bank uint16, program uint8, bytes ...byte) ProgramChange {
	if channel > 15 {
		panic(fmt.Sprintf("invalid channel (%v)", channel))
	}

	return ProgramChange{
		event: event{
			tick:    tick,
			delta:   lib.Delta(delta),
			bytes:   bytes,
			tag:     lib.TagProgramChange,
			Status:  or(0xc0, channel),
			Channel: channel,
		},
		Bank:    bank,
		Program: program,
	}
}

func UnmarshalProgramChange(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (*ProgramChange, error) {
	if status&0xf0 != 0xc0 {
		return nil, fmt.Errorf("Invalid %v status (%v): expected 'Cx'", lib.TagProgramChange, status)
	}

	if len(data) < 1 {
		return nil, fmt.Errorf("Invalid %v data (%v): expected note and velocity", lib.TagProgramChange, data)
	}

	var channel = lib.Channel(status & 0x0f)
	var bank uint16
	var program = data[0]

	if ctx != nil {
		bank = ctx.ProgramBank[uint8(channel)]
	}

	event := MakeProgramChange(tick, delta, channel, bank, program, bytes...)

	return &event, nil
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
	p.tag = lib.TagProgramChange
	p.Status = 0xc0

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)ProgramChange\s+channel:([0-9]+)\s+bank:([0-9]+),\s*program:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 5 {
		return fmt.Errorf("invalid ProgramChange event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if channel, err := lib.ParseChannel(match[2]); err != nil {
		return err
	} else if bank, err := strconv.ParseUint(match[3], 10, 16); err != nil {
		return err
	} else if program, err := strconv.ParseUint(match[4], 10, 8); err != nil {
		return err
	} else {
		p.delta = lib.Delta(delta)
		p.Status = or(p.Status, channel)
		p.Channel = channel
		p.Bank = uint16(bank)
		p.Program = uint8(program)
	}

	return nil
}
